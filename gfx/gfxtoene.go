package gfx
// Autor: St. Schmidt
// Datum: 08.05.2019
// Zweck: Töne/Noten spielen können

import "math"

const (             // Konstanten für die Signalform der gespielten Töne
Sinusform uint8 = iota
Rechteckform
Dreieckform
Sägezahnform
)

var r uint32= 44100 // Abtastrate: 11025 oder 22050 oder 44100 - Standard 44100 Hz
var b uint8 = 2     // Auflösung:    1: 8 Bit ; 2: 16 Bit      - Standard 16 Bit
var k uint8 = 2     // Kanalanzahl:  1: mono  ; 2: stereo      - Standard stereo

var s = Rechteckform // aktuelle Signalform:                 - Standard Rechteck

var anschlagzeit  float64 = 0.002 // Standard auf 2 ms gesetzt
var abschwellzeit float64 = 0.750 // Standard auf 750 ms gesetzt
var haltepegel    float64 = 0     // Standard auf 0 % gesetzt ; 1 = 100 %
var ausklingzeit  float64 = 0.006 // Standard auf 6 ms gesetzt
var pulsweite     float64 = 0.375 // nur wichtig für Rechteck-Signale, hier: Prozentsatz HIGH (Pulsweite)

var tempo uint8 = 120 // "Schläge" = Viertelnoten pro Minute"
var tVollnote uint16 = 4 * uint16(((60 * 1000 + uint32(tempo)/2 )/uint32(tempo))) // Zeit einer vollen Note in ms

var f  = map[string]float64{ // Frequenzen der Noten - 7. Oktave
	"C" :2093.00, "C#":2217.46,
	"D" :2349.32, "D#":2489.02,
	"E" :2637.02,
	"F" :2793.83, "F#":2959.96,
	"G" :3135.96, "G#":3322.44,
	"A" :3520.00, "A#":3729.31,
	"H" :3951.07}

// Erg.: Das aktuelle Tempo ist geliefert, d.h. die Anzahl der Viertelnoten pro Minute.
func GibNotenTempo () uint8 { return tempo }

// Vor.: 30 <= tempo <= 240 
// Eff.: Das Tempo ist auf den Wert t gesetzt.
func SetzeNotenTempo (t uint8) {
	if t >= 30  && t <= 240 {
		tempo = t
		tVollnote = 4 * uint16(((60 * 1000 + uint32(tempo)/2 )/uint32(tempo)))
	}
}
		
// Vor.: -
// Erg.: Geliefert sind:
//       Abtastrate der WAV-Daten in Hz, z.B. 44100,
//       Auflösung der Klänge (1: 8 Bit; 2: 16 Bit),
//       die Anzahl der Kanäle (1: mono, 2:stereo),
//       die Signalform (0: Sinus, 1: Rechteck, 2:Dreieck, 3: Sägezahn) und
//       die Pulsweite HIGH bei Rechteckform als Prozentsatz zw. 0 und 1.
func GibKlangparameter () (uint32,uint8,uint8,uint8,float64) {
	return r,b,k,s,pulsweite
}
		
// Vor.: rate ist die Abtastrate, z.B. 11025, 22050 oder 44100.
//       auflösung ist 1 für 8 Bit oder 2 für 16 Bit.
//       kanaele ist 1 für mono oder 2 für stereo.
//       signal gibt die Signalform an: 0: Sinus, 1: Rechteck, 2:Dreieck, 3: Sägezahn
//       p ist die Pulsweite für Rechtecksignale und gibt den Prozentsatz (0<=p<=1) für den HIGH-Teil an.
// Eff.: Die klangparameter sind auf die angegebenen Werte gesetzt.
func SetzeKlangparameter(rate uint32, aufloesung,kanaele,signal uint8, p float64) {
	r = rate
	b = aufloesung
	k = kanaele
	s = signal
	pulsweite = p
}
		
// Vor.: -
// Ergebnis: Anschlagzeit, Abschwellzeit, Haltepegel und Ausklingzeit sind geliefert.
func GibHuellkurve () (float64,float64,float64,float64) {
	return anschlagzeit,abschwellzeit,haltepegel,ausklingzeit
}

// Vor.: a ist die Anschlagzeit in s mit 0 <= a <= 1,
//       d ist die Abschwellzeit in s mit 0<= d <= 5,
//       s ist der Haltepegel in Prozent vom Maximum mit 0<= s <= 1.0,
//       r ist die Ausklingzeit in s mit 0< =r <= 5.
// Eff.: Für die Hüllkurve zukünftig zu spielender Töne bzw. Noten sind
//       die Parameter gesetzt.
func SetzeHuellkurve (a,d,s,r float64) {
	if 0<=a && a<= 1 && 0<=d && d<=5 && 0<=s && s<=1 && 0<=r && r<=5 {
		anschlagzeit  = a
		abschwellzeit = d
		haltepegel    = s
		ausklingzeit  = r
	}
}

// INTERN
// Erg.: Ein Slice bestehend aus 4 Bytes ist geliefert, die den Wert x 
//       darstellen. Das erste Byte ist das LSB.
func uint32toSlice (x uint32) []byte {
	var erg []byte = make ([]byte,4)
	for i:=0; i<4;i++ {
		erg[i] = byte(x % 256)
		x = x / 256
	}
	return erg
}

// INTERN
// Vor.: t ist der echte Zeitpunkt innerhalb des Tons.
// Erg.: Der Maximalausschlag zwischen 0 und 1 ist gemäß
//       der aktuell festgelegten Hüllkurve geliefert. 
func amplitude (t,tges float64) float64 {
	switch {
		case t <= anschlagzeit:
		return t/anschlagzeit
		case t > tges-ausklingzeit:
		return haltepegel-(t-(tges-ausklingzeit))*haltepegel/ausklingzeit
		default: //Abschwell- und Haltepegelzeit
		return haltepegel+(1-haltepegel)*math.Pow(2,-(t-anschlagzeit)*6/abschwellzeit)
	}
}

// INTERN
// Vor.: tges gibt die Tondauer in Millisekunden an.
//       f ist die Tonfrequenz in Hertz.
// Erg.: Ein Byte-Slice ist geliefert, dass der entsprechenden WAV-Datei entspricht.
func ton (tges uint16, f float64) []byte {
	var laenge uint32           = r*uint32(tges)*uint32(b*k)/1000
	var dateigrößeMinus8 uint32 = laenge + 44 - 8 
	var bytes []byte            = make ([]byte,44 + laenge)
	var w float64
	
	// "Dateikopf gemäß RIFF-WAVE-Format
	copy (bytes,"RIFF")
	copy (bytes[4:], uint32toSlice(dateigrößeMinus8)) //DATEIGRÖSSE - 8 ----------------------------------------------------
	copy (bytes[8:],"WAVEfmt ")
	bytes[16] = 16 // Die Größe des fmt-Abschnitts ist 16 Bytes (uint32)
	bytes[17] = 0
	bytes[18] = 0
	bytes[19] = 0
	bytes[20] = 1 // Das verwendete Format: 01 = PCM (uint16)
	bytes[21] = 0
	bytes[22] = k // Wir verwenden k Kanal (1:mono; 2:stereo).
	bytes[23] = 0
	copy (bytes[24:],uint32toSlice(r)) // Eintrag der Abtastrate
	copy (bytes[28:],uint32toSlice(r * uint32(b*k))) // Übertragungsbandbreite (Bytes pro Sekunde): rate*b*k
	bytes[32] = k   // uint16 - 1: mono ; 2: stereo
	bytes[33] = 0
	bytes[34] = b*8 // uint16 - Auflösung: 8 oder 16 Bit
	bytes[35] = 0
	copy(bytes[36:],"data")
	copy(bytes[40:], uint32toSlice(laenge))//DATEIGRÖSSE - 44----------------------------------------------------
	// Es folgen die Daten - ein Frame = b*k Byte
	// Es sind rate Frames pro Sekunde
	for i:=uint32(0);i<laenge-uint32(b*k)+1;i=i+uint32(b*k) {
		t:= float64(i)/float64(r*uint32(b*k)) // echter Zeitpunkt 
		t2:= t-float64(uint64(t*f))/f         // Zeitpunkt innerhalb der aktuellen Schwingung 
		switch s { // nach Signalform
			case Sinusform:
			w= math.Sin(2*math.Pi*f*t2) // float64 aus [-1;1]
			case Rechteckform:
			if t2 <= pulsweite/f {
				w = 1
			} else {
				w = -1
			}
			case Dreieckform:
			if t2 <= 1/(2*f) {
				w = -1 + 4*f*t2
			} else {
				w = 1 - 4 * f * (t2-1/(2*f))
			}
			case Sägezahnform:
			w=-1+2*f*t2
			default:
			panic ("unbekannte Signalform!!")
		}
		w = amplitude(t,float64(tges)/1000) * w  // Einarbeiten der Hüllkurve
		switch b {
			case 1:
			bytes[44+i] = uint8 (128 + w * 127)
			if k == 2 { 
				bytes[45+i] = bytes[44+i]
			}
			case 2:
			bytes[44+i] = byte(uint16(w*32767) % 256)
			bytes[45+i] = byte(uint16(w*32767) / 256)
			if k == 2 {
				bytes[46+i] = bytes[44+i]
				bytes[47+i] = bytes[45+i]
			}
		}
	}	
	return bytes
}

// Vor.: Das erste Zeichen von tonname ist eine Ziffer von 0 bis 9 und gibt die Oktave an.
//       Erlaubte weitere Zeichen für den Notennamen sind "C","D","E","F","G","A","H","C#","D#","F#","G#","A#".
//       0 < laenge <= 1;  laenge 1: volle Note; 1.0/2: halbe Note, ..., 1.0/16: sechzehntel Note
//       0.0<=wartedauer<=2.0; Die Wartedauer gibt die Dauer in Notenlänge an, nach der nach dem Anspielen der
//       Note im Programmablauf fortgefahren wird. 0: keine Wartedauer; 1.0/2: Dauer einer halben Note, ...  
// Eff.: Der Ton wird gerade gespielt bzw. ist gespielt. Je nach Wartedauer wurde die Fortsetzung des Programms
//       verzögert.
//       Der voreingestellte Standard ist aus 'GibHuellkurve ()' und 'GibKlangParameter()' ersichtlich.
//       Die Einstellungen mit 'SetzeHuellkurve' und 'SetzeKlangparameter' haben Einfluss auf den "Ton".
func SpieleNote (tonname string, laenge float64, wartedauer float64) {
	var o uint8 = byte(tonname[0])-48
	var freq float64 =f[tonname[1:]]
	for i:=uint8(7);i>o;i--{ freq = freq / 2 }
	for i:=uint8(7);i<o;i++{ freq = freq * 2 }  
	bytes:= ton(uint16(float64(tVollnote)*laenge),freq)
	spieleRAMWAV(bytes,uint32(wartedauer*float64(tVollnote)))
}

