package gfx

// Autor: Stefan Schmidt (Kontakt: St.Schmidt@online.de)
// Datum: 07.03.2016 ; letzte Änderung: 01.09.2019
// Zweck: - Grafik- und Soundausgabe und Eingabe per Tastatur und Maus
//          mit Go unter Windows und unter Linux
//        - 01.09.2019  Spezifikationsfehler korrigiert
//        - 08.05.2019  neue Funktionen, um (Klavier-) Noten spielen zu können
//                      inkl. Hüllkurven und klanganpassungen
//        - 07.03.2019  Rechtschreibkorrekturen in der Spezifikation
//        - 03.03.2019: Die Funktion 'SetzeFont' liefert nun einen Rückgabewert,
//                      der den Erfolg/Misserfolg angibt.
//        - 07.10.2017: neue Funktion 'Tastaturzeichen'
//        - 07.10.2017: 'Bug' in Funktion 'Cls()' entfernt - KEIN FLACKERN MEHR
//                       bei 'double-buffering' mit UpdateAus() und UpdateAn()

/*
#cgo LDFLAGS: -lSDL -lSDL_gfx -lSDL_ttf
#include <SDL/SDL.h>
#include <SDL/SDL_ttf.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <SDL/SDL_gfxPrimitives.h>

// Structure for loaded sounds.
typedef struct sound_s {
    Uint8 *samples;		// raw PCM sample data
    Uint32 length;		// size of sound data in bytes
} sound_t, *sound_p;

// Structure for a currently playing sound.
typedef struct playing_s {
    int active;                 // 1 if this sound should be played
    sound_p sound;              // sound data to play
    Uint32 position;            // current position in the sound buffer
} playing_t, *playing_p;

// Array for all active sound effects.
#define MAX_PLAYING_SOUNDS      10
playing_t playing[MAX_PLAYING_SOUNDS];

// The higher this is, the louder each currently playing sound will be.
// However, high values may cause distortion if too many sounds are
// playing. Experiment with this.
#define VOLUME_PER_SOUND        SDL_MIX_MAXVOLUME / 2

static SDL_Surface *screen;
static SDL_Surface *archiv;
static SDL_Surface *clipboard = NULL;
static Uint8 updateOn = 1;
static Uint8 red,green,blue, alpha;
static SDL_Event event;
static Uint8 gedrueckt;
static Uint16 taste,tiefe;
static Uint8 tasteLesen = 0;
static Uint8 tastaturpuffer = 0;
static Uint32 t_puffer[256];
static Uint8 t_pufferkopf;
static Uint8 t_pufferende;
static Uint16 mausX, mausY;
static Uint8 mausLesen = 0;
static Uint8 mausTaste;
static Uint8 mauspuffer = 0;
static Uint32 m_puffer[256];
static Uint8 m_pufferkopf;
static Uint8 m_pufferende;
static Uint8 fensteroffen = 0;
static Uint8 fensterzu = 1;
static char aktFont[256];
static int aktFontSize;
static SDL_AudioSpec desired, obtained; // Audio format specifications.
static sound_t s[10];                   // Our loaded sounds and their formats.

//------------------------------------------------------------------
// This function is called by SDL whenever the sound card
// needs more samples to play. It might be called from a
// separate thread, so we should be careful what we touch.
static void AudioCallback(void *user_data, Uint8 *audio, int length)
{
    int i;
    // Avoid compiler warning.
    user_data += 0;
    // Clear the audio buffer so we can mix samples into it.
    memset(audio, 0, length);
    // Mix in each sound.
    for (i = 0; i < MAX_PLAYING_SOUNDS; i++) {
	  if (playing[i].active) {
	    Uint8 *sound_buf;
	    Uint32 sound_len;
	    // Locate this sound's current buffer position.
	    sound_buf = playing[i].sound->samples;
	    sound_buf += playing[i].position;
	    // Determine the number of samples to mix.
	    if ((playing[i].position + length) > playing[i].sound->length) {
		sound_len = playing[i].sound->length - playing[i].position;
	    } else {
		sound_len = length;
	    }
	    // Mix this sound into the stream.
	    SDL_MixAudio(audio, sound_buf, sound_len, VOLUME_PER_SOUND);
	    // Update the sound buffer's position.
	    playing[i].position += length;
	    // Have we reached the end of the sound?
	    if (playing[i].position >= playing[i].sound->length) {
	    free(s[i].samples);      //zugehörigen Soundstruktur-Samplespeicher wieder freigeben
		playing[i].active = 0;	 // und anschließend als inaktiv markieren
	    }
	  }
    }
}
//----------------------------------------------------------------
// This function loads a sound with SDL_LoadWAV and converts
// it to the specified sample format. Returns 0 on success
// and 1 on failure.
static int LoadAndConvertSound(char *filename, SDL_AudioSpec *spec,
			sound_p sound)
{
    SDL_AudioCVT cvt;           // audio format conversion structure
    SDL_AudioSpec loaded;       // format of the loaded data
    Uint8 *new_buf;
    // Load the WAV file in its original sample format.
    if (SDL_LoadWAV(filename,
		    &loaded, &sound->samples, &sound->length) == NULL) {
	printf("Unable to load sound: %s\n", SDL_GetError());
	return 1;
    }
    // Build a conversion structure for converting the samples.
    // This structure contains the data SDL needs to quickly
    // convert between sample formats.
    if (SDL_BuildAudioCVT(&cvt, loaded.format,
			  loaded.channels,
			  loaded.freq,
			  spec->format, spec->channels, spec->freq) < 0) {
	// printf("Unable to convert sound: %s\n", SDL_GetError());
	return 1;
    }
    // Since converting PCM samples can result in more data
    //   (for instance, converting 8-bit mono to 16-bit stereo),
    //   we need to allocate a new buffer for the converted data.
    //   Fortunately SDL_BuildAudioCVT supplied the necessary
    //   information.
    cvt.len = sound->length;
    new_buf = (Uint8 *) malloc(cvt.len * cvt.len_mult);
    if (new_buf == NULL) {
	//printf("Memory allocation failed.\n");
	SDL_FreeWAV(sound->samples);
	return 1;
    }
    // Copy the sound samples into the new buffer.
    memcpy(new_buf, sound->samples, sound->length);
    // Perform the conversion on the new buffer.
    cvt.buf = new_buf;
    if (SDL_ConvertAudio(&cvt) < 0) {
	//printf("Audio conversion error: %s\n", SDL_GetError());
	free(new_buf);
	SDL_FreeWAV(sound->samples);
	return 1;
    }
    // Swap the converted data for the original.
    SDL_FreeWAV(sound->samples);
    sound->samples = new_buf;
    sound->length = sound->length * cvt.len_mult;
    // Success!
    //printf("'%s' was loaded and converted successfully.\n", filename);
    return 0;
}
//----------------------------------------------------------------
// Diese Funktion übernimmt eine Bytefolge aus dem RAM ab der Adresse addr
// mit der Länge laenge, die dem Inhalt einer WAV-DAtei entspricht und konvertiert
// sie , damit es abgespielt werden kann. Die Funktion liefert 0 bei Erfolg
// 1 bei Misserfolg.
static int LadeUndKonvertiereRAMWAV(const void* addr, int laenge, SDL_AudioSpec *spec,
			sound_p sound)
{
    SDL_AudioCVT cvt;           // audio format conversion structure
    SDL_AudioSpec loaded;       // format of the loaded data
    Uint8 *new_buf;
    // Lade 'RAMWAV' im Originalformat:
    if (SDL_LoadWAV_RW(SDL_RWFromConstMem(addr,laenge),0,
		    &loaded, &sound->samples, &sound->length) == NULL) {
	printf("Unable to load sound: %s\n", SDL_GetError());
	return 1;
    }
    // Build a conversion structure for converting the samples.
    // This structure contains the data SDL needs to quickly
    // convert between sample formats.
    if (SDL_BuildAudioCVT(&cvt, loaded.format,
			  loaded.channels,
			  loaded.freq,
			  spec->format, spec->channels, spec->freq) < 0) {
	// printf("Unable to convert sound: %s\n", SDL_GetError());
	return 1;
    }
    // Since converting PCM samples can result in more data
    //   (for instance, converting 8-bit mono to 16-bit stereo),
    //   we need to allocate a new buffer for the converted data.
    //   Fortunately SDL_BuildAudioCVT supplied the necessary
    //   information.
    cvt.len = sound->length;
    new_buf = (Uint8 *) malloc(cvt.len * cvt.len_mult);
    if (new_buf == NULL) {
	//printf("Memory allocation failed.\n");
	SDL_FreeWAV(sound->samples);
	return 1;
    }
    // Copy the sound samples into the new buffer.
    memcpy(new_buf, sound->samples, sound->length);
    // Perform the conversion on the new buffer.
    cvt.buf = new_buf;
    if (SDL_ConvertAudio(&cvt) < 0) {
	//printf("Audio conversion error: %s\n", SDL_GetError());
	free(new_buf);
	SDL_FreeWAV(sound->samples);
	return 1;
    }
    // Swap the converted data for the original.
    SDL_FreeWAV(sound->samples);
    sound->samples = new_buf;
    sound->length = sound->length * cvt.len_mult;
    // Success!
    //printf("'%s' was loaded and converted successfully.\n", filename);
    return 0;
}
//-----------------------------------------------------------------
static int LoadAndPlaySound (char *filename)
{
	int i;
	//Finde einen freien Index (Bereich 0 <= index < MAX_PLAYING_SOUND
	for (i = 0; i < MAX_PLAYING_SOUNDS; i++) {
	if (playing[i].active == 0)
	    break;
    }
    if (i == MAX_PLAYING_SOUNDS)
	return 1; //Fehler: Es werden schon die max. Anzahl an Dateien abgespielt.

	//Lade und konvertiere den Sound in die entsprechende Soundstruktur
	if (LoadAndConvertSound(filename, &obtained, &s[i]) != 0) {
	  return 2; //Laden fehlgeschlagen!
    }
    //Abspielen starten
    // The 'playing' structures are accessed by the audio callback,
    // so we should obtain a lock before we access them.
    SDL_LockAudio();
    playing[i].active = 1;
    playing[i].sound = &s[i];
    playing[i].position = 0;
    SDL_UnlockAudio();
    return 0;
}
//-----------------------------------------------------------------
static int LadeUndSpieleNote (const void* addr, int laenge)
{
	int i;
	//Finde einen freien Index (Bereich 0 <= index < MAX_PLAYING_SOUND
	for (i = 0; i < MAX_PLAYING_SOUNDS; i++) {
	if (playing[i].active == 0)
	    break;
    }
    if (i == MAX_PLAYING_SOUNDS)
	return 1; //Fehler: Es werden schon die max. Anzahl an Dateien abgespielt.

	//Lade und konvertiere den Sound in die entsprechende Soundstruktur
	if (LadeUndKonvertiereRAMWAV(addr, laenge, &obtained, &s[i]) != 0) {
	  return 2; //Laden fehlgeschlagen!
    }
    //Abspielen starten
    // The 'playing' structures are accessed by the audio callback,
    // so we should obtain a lock before we access them.
    SDL_LockAudio();
    playing[i].active = 1;
    playing[i].sound = &s[i];
    playing[i].position = 0;
    SDL_UnlockAudio();
    return 0;
}
//------------------------------------------------------------------
static int setFont (char *fontfile, int groesse) {
	strcpy (aktFont,fontfile);
	aktFontSize = groesse;
	TTF_Font *font = TTF_OpenFont(aktFont, aktFontSize);
	if (!font) {
	  //printf("TTF_OpenFont: %s\n", TTF_GetError());
	  return 1;
    }
    TTF_CloseFont(font);
	return 0;
}
//------------------------------------------------------------------
static char *getFont () {
	return aktFont;
}
//------------------------------------------------------------------
static int write (Sint16 x, Sint16 y, char *text) {
	TTF_Font *font = TTF_OpenFont(aktFont, aktFontSize);
	if (!font) {
	  //printf("TTF_OpenFont: %s\n", TTF_GetError());
	  return 1;
    }
	SDL_Color clrFg = {red,green,blue,alpha};
	SDL_Surface *sText = TTF_RenderUTF8_Solid(font,text,clrFg);
	SDL_Rect rcDest = {x,y,0,0};
	SDL_BlitSurface(sText,NULL, screen,&rcDest);
	SDL_FreeSurface(sText);
	if (updateOn)
	  SDL_UpdateRect(screen,0,0,0,0);
	TTF_CloseFont(font);
	return 0;
}
//------------------------------------------------------------------
static void clearscreen () {
  SDL_FillRect(screen, NULL, SDL_MapRGB(screen->format, red, green, blue));
  if (updateOn)
    SDL_UpdateRect (screen,0,0,0,0);
}
//------------------------------------------------------------
static int GrafikfensterAn (Uint16 breite, Uint16 hoehe)
{
    if ( fensteroffen == 1) return 1;  //Es kann nur ein Grafikfenster geben!

	//1. SDL muss initialisiert werden.
	if (SDL_Init(SDL_INIT_VIDEO | SDL_INIT_AUDIO) != 0) {
		//printf ("Kann SDL nicht initialisieren: %s\n", SDL_GetError ());
		return 1;
	}
	//2. Bekanntmachung: Diese Funktion soll mit dem Programmende aufgerufen werden.
	// atexit (SDL_Quit);
	//3. Bildschirm: Hier kann man auch SDL_DOUBLEBUF sagen!
	screen = SDL_SetVideoMode (breite,hoehe, 32, SDL_DOUBLEBUF); //SDL_FULLSCREEN);
	if (screen == NULL) {
		//printf ("Bildschirm-Modus nicht setzbar: %s\n",SDL_GetError ());
		return 1;
	}
	SDL_WM_SetCaption( "LWB FU-Berlin: GO-Grafikfenster", 0 );

	TTF_Init();

	red = 255;
	green = 255;
	blue = 255;
	alpha = 255;
	clearscreen ();
	red   = 0;
	green = 0;
	blue  = 0;

	//Archiv-Surface erstellen
	archiv = SDL_ConvertSurface (screen, screen->format, SDL_HWSURFACE);
	if (archiv == NULL) {
		//printf ("Archiv-Surface konnte nicht erzeugt werden!\n");
		return 1;
	}

    // Open the audio device. The sound driver will try to give us
    // the requested format, but it might not succeed. The 'obtained'
    // structure will be filled in with the actual format data.
    desired.freq = 44100;	// desired output sample rate
    desired.format = AUDIO_S16;	// request signed 16-bit samples
    desired.samples = 8192;	// this is more or less discretionary
    desired.channels = 2;	// ask for stereo
    desired.callback = AudioCallback;
    desired.userdata = NULL;	// we don't need this
    if (SDL_OpenAudio(&desired, &obtained) < 0) {
    	//printf("Unable to open audio device: %s\n", SDL_GetError());
	    return 1;
    }
    // Initialisiere die Liste der möglichen Sounds (keiner aktiv zu Beginn)
    int i;
    for (i = 0; i < MAX_PLAYING_SOUNDS; i++) {
	playing[i].active = 0;
    }

    // SDL's audio is initially paused. Start it.
    SDL_PauseAudio(0);

	fensteroffen = 1;
	fensterzu = 0;

	//Jetzt kommt die Event-Loop

	while (SDL_WaitEvent(&event) != 0 && fensteroffen == 1) {
		switch (event.type) {
			case SDL_KEYDOWN:
				if (tasteLesen)
				{
					gedrueckt = 1;                //Taste ist gerade heruntergedrückt.
					taste = event.key.keysym.sym;  //Das ist der Code der Taste auf der Tastatur.
					tiefe = event.key.keysym.mod;  //Gleichzeitig Steuerungstaste(n) gedrückt??
					//printf("%i,%i,%i\n",taste, gedrueckt, tiefe);
					tasteLesen = 0;
				}
				if (tastaturpuffer)
				{
					if (t_pufferende + 1 != t_pufferkopf)
					{
						t_puffer[t_pufferende] = ((Uint32) event.key.keysym.sym)*256*256 + (Uint32) 256*256*256*128 + ((Uint32) event.key.keysym.mod);
						t_pufferende++; //Umschlag auf 0 automatisch, da Uint8
					}
				}
				break;
			case SDL_KEYUP:
				if (tasteLesen)
				{
					gedrueckt = 0; //Taste wurde gerade losgelassen.
					taste = event.key.keysym.sym;
					tiefe = event.key.keysym.mod;  //Gleichzeitig Steuerungstaste(n) gedrückt??
					//printf("%i,%i,%i\n",taste, gedrueckt, tiefe);
					tasteLesen = 0;
				}
				if (tastaturpuffer)
				{
					if (t_pufferende + 1 != t_pufferkopf)
					{
						t_puffer[t_pufferende] = ((Uint32) event.key.keysym.sym)*256*256 + ((Uint32) event.key.keysym.mod);
						t_pufferende++; //Umschlag auf 0 automatisch, da Uint8
					}
				}
				break;
			case SDL_MOUSEMOTION:
				if (mausLesen)
				{   //BEi MOUSEMOTION GIBT ES NUR 3 MÖGLICHKEITEN FÜR EINE GEDRÜCKT-GEHALTENE TASTE: 1,2 oder 3
					// Dummerweise ist bei 3 der Tastenwert 4, daher Korrektur:
					mausTaste = (Uint8) event.button.button;
					if (mausTaste == 4)
						mausTaste--;
					mausX     = (Uint16) event.motion.x;
					mausY     = (Uint16) event.motion.y;
					mausLesen = 0;
				}
				if (mauspuffer)
				{
					mausTaste = (Uint8) event.button.button;
					if (mausTaste == 4)
						mausTaste--;
					mausX     = (Uint16) event.motion.x;
					mausY     = (Uint16) event.motion.y;
					if (m_pufferende + 1 != m_pufferkopf)
					{
						m_puffer[m_pufferende] = ((Uint32) mausTaste)*256*256*256 + (((Uint32) mausX) <<12) + (Uint32) mausY;
						m_pufferende++;
					}
				}
				break;
			case SDL_MOUSEBUTTONDOWN:
				if (mausLesen)
				{
					mausTaste = (Uint8) event.button.button + 128; //+128: "pressed"
					mausX = (Uint16) event.motion.x;
					mausY = (Uint16) event.motion.y;
					mausLesen = 0;
				}
				if (mauspuffer)
				{
					mausTaste = (Uint8) event.button.button + 128; //+128: "pressed"
					mausX = (Uint16) event.motion.x;
					mausY = (Uint16) event.motion.y;
					if (m_pufferende + 1 != m_pufferkopf)
					{
						m_puffer[m_pufferende] = ((Uint32) mausTaste)*256*256*256 + (((Uint32) mausX) <<12) + (Uint32) mausY;
						m_pufferende++;
					}
				}
				break;
			case SDL_MOUSEBUTTONUP:
				if (mausLesen)
				{
					mausTaste = (Uint8) event.button.button + 64; //+64: "released"
					mausX = (Uint16) event.motion.x;
					mausY = (Uint16) event.motion.y;
					mausLesen = 0;
				}
				if (mauspuffer)
				{
					mausTaste = (Uint8) event.button.button + 64; //+64: "released"
					mausX = (Uint16) event.motion.x;
					mausY = (Uint16) event.motion.y;
					if (m_pufferende + 1 != m_pufferkopf)
					{
						m_puffer[m_pufferende] = ((Uint32) mausTaste)*256*256*256 + (((Uint32) mausX) <<12) + (Uint32) mausY;
						m_pufferende++;
					}
				}
				break;
			case SDL_QUIT:
				//printf("Das Grafikfenster wurde geschlossen. Bye.\n");
				exit(0);
		}

	}

	// Die event-Loop wurde beendet, also wird nun das Fenster geschlossen!
	TTF_Quit ();
    // Pause and lock the sound system so we can safely delete our sound data.
    SDL_PauseAudio(1);
    SDL_LockAudio();
    // Free our sounds before we exit, just to be safe.
    for (i=0; i < MAX_PLAYING_SOUNDS;i++) {
		if (playing[i].active ==1) {
			free(s[i].samples);
		}
	}
    // At this point the output is paused and we know for certain that the
    // callback is not active, so we can safely unlock the audio system.
    SDL_UnlockAudio();
	SDL_CloseAudio();
	SDL_Quit ();
	fensterzu = 1;
	return 0;
}
//------------------------------------------------------------------
static Uint8 FensterOffen ()
{
  return fensteroffen;
}
//------------------------------------------------------------------
static Uint8 FensterZu ()
{
  return fensterzu;
}
//-------------------------------------------------------------------
static void GrafikfensterAus ()
{
  fensteroffen = 0;
}
//-------------------------------------------------------------------
static void updateAus ()
{
	updateOn = 0;
}
//-------------------------------------------------------------------
static void updateAn ()
{
	updateOn = 1;
	SDL_Flip (screen);
}
//-------------------------------------------------------------------
static void zeichnePunkt (Sint16 x, Sint16 y)
{
  pixelRGBA (screen, x, y ,red, green,blue,alpha);
  if (updateOn)
	SDL_UpdateRect (screen, x, y, 1, 1);
}
//--------------------------------------------------------------
static Uint32 gibPixel(Sint16 x, Sint16 y)
{
    int bpp = screen->format->BytesPerPixel;
    // Here p is the address to the pixel we want to retrieve
    Uint8 *p = (Uint8 *)screen->pixels + y * screen->pitch + x * bpp;

    switch(bpp) {
    case 1:
        return *p;
        break;
    case 2:
        return *(Uint16 *)p;
        break;
    case 3:
        if(SDL_BYTEORDER == SDL_BIG_ENDIAN)
            return p[0] << 16 | p[1] << 8 | p[2];
        else
            return p[0] | p[1] << 8 | p[2] << 16;
        break;
    case 4:
        return *(Uint32 *)p;
        break;
    default:
        return 0;       // shouldn't happen, but avoids warnings
    }
}
//--------------------------------------------------------------
static void zeichneKreis (Sint16 x, Sint16 y, Sint16 r, Uint8 full)
{
	if (full)
		filledCircleRGBA(screen,x,y,r,red,green,blue,alpha);
	else
		circleRGBA (screen, x,y,r,red, green, blue,alpha);
	if (updateOn)
		SDL_UpdateRect (screen,x-r,y-r,2*r+1,2*r+1);
}
//---------------------------------------------------------------
static void zeichneEllipse (Sint16 x, Sint16 y, Sint16 rx, Sint16 ry, Uint8 filled)
{
	if (filled)
		filledEllipseRGBA (screen, x, y, rx, ry, red, green, blue, alpha);
	else
		ellipseRGBA (screen, x, y, rx, ry, red, green,blue, alpha);
	if (updateOn)
		SDL_UpdateRect (screen, x-rx, y-ry,2*rx+1,2*ry+1);
}
//---------------------------------------------------------------
static void stiftfarbe (Uint8 r, Uint8 g, Uint8 b)
{
    red = r;
    green = g;
    blue = b;
}
//---------------------------------------------------------------
static void zeichneStrecke (Sint16 x1, Sint16 y1, Sint16 x2, Sint16 y2)
{
	int upx,upy,breite,hoehe;

	lineRGBA (screen, x1,y1,x2,y2, red, green, blue, alpha);
	if (x1 <= x2)
	{
		upx    = x1;
		breite = x2 - x1 + 1;
	}
	else
	{
		upx    = x2;
		breite = x1 - x2 + 1;
	}
	if (y1 <= y2)
	{
		upy   = y1;
		hoehe = y2 - y1 + 1;
	}
	else
	{
		upy   = y2;
		hoehe = y1 - y2 + 1;
	}
	if (updateOn)
		SDL_UpdateRect (screen,upx,upy,breite,hoehe);
}
//--------------------------------------------------------------
static void rechteck (Sint16 x1, Sint16 y1, Sint16 b, Sint16 h, Uint8 filled)
{
	if (filled)
		boxRGBA (screen, x1, y1 ,x1+b-1, y1+h-1, red, green, blue, alpha);
	else
		rectangleRGBA (screen, x1, y1 , x1+b-1, y1+h-1, red, green,blue,alpha);
	if (updateOn)
		SDL_UpdateRect (screen, x1, y1, b, h);
}
//--------------------------------------------------------------
static void kreissektor (Sint16 x, Sint16 y, Sint16 r, Sint16 w1, Sint16 w2, Uint8 filled)
{
	if (filled)
		filledPieRGBA (screen, x, y , r, w1, w2, red, green, blue, alpha);
	else
		pieRGBA (screen, x, y , r, w1, w2, red, green, blue, alpha);
	if (updateOn)
		SDL_UpdateRect (screen, x-r, y-r, 2*r+1, 2*r+1);
}
//---------------------------------------------------------------
Sint16 minimum (Sint16 x, Sint16 y, Sint16 z)
{
  if ((x <= y) && (x <=z))
    return x;
  else if ((y<=x) && (y<=z))
    return y;
  else
    return z;
}
//-------------------------------------------------------------
Sint16 maximum (Sint16 x, Sint16 y, Sint16 z)
{
  if ((x >= y) && (x >=z))
    return x;
  else if ((y>=x) && (y>=z))
    return y;
  else
    return z;
}
//---------------------------------------------------------------
static void dreieck (Sint16 x1, Sint16 y1, Sint16 x2, Sint16 y2, Sint16 x3, Sint16 y3, Uint8 filled)
{
	int upx,upy,breite,hoehe;

	upx = minimum (x1, x2, x3);
	upy = minimum (y1, y2, y3);
	breite = maximum (x1, x2, x3) - upx + 1;
	hoehe  = maximum (y1, y2, y3) - upy + 1;
	if (filled)
		filledTrigonRGBA(screen, x1,y1,x2,y2,x3,y3,red,green,blue,alpha);
	else
		trigonRGBA (screen, x1,y1,x2,y2,x3,y3,red,green,blue,alpha);
	if (updateOn)
		SDL_UpdateRect (screen, upx,upy,breite,hoehe);
}
//----------------------------------------------------------------
static void ladeBild (Sint16 x, Sint16 y, char *cs)
{
	SDL_Surface *image;
	SDL_Rect src, dest;

	image = SDL_LoadBMP(cs);
	//printf ("Dateiname: %s\n",cs);
	if (image == NULL) {
		//printf("Bild konnte nicht geladen werden!\n");
		return;
	}
	src.x = 0;
	src.y = 0;
	src.w = image->w;
	src.h = image->h;

	dest.x = x;
	dest.y = y;
	dest.w = image->w;
	dest.h = image->h;

	SDL_BlitSurface(image, &src, screen, &dest);
	SDL_FreeSurface (image);
	if (updateOn)
		SDL_UpdateRect(screen, x, y, dest.w, dest.h);
}
//---------------------------------------------------------
static void schreibe (Sint16 x, Sint16 y, char *cs)
{
	gfxPrimitivesSetFont(NULL, 0 ,0);
	stringRGBA (screen, x,y,cs,red, green, blue, alpha);
	if (updateOn)
		SDL_UpdateRect (screen,0,0,0,0);
}
//---------------------------------------------------------
static void ladeBildInsClipboard (char *cs)
{
	SDL_Surface *image;

	image = SDL_LoadBMP(cs);
	//printf ("Dateiname: %s\n",cs);
	if (image == NULL) {
		// printf("Bild konnte nicht geladen werden!\n");
		return;
	}
	SDL_FreeSurface (clipboard); //altes Clipboard freigeben
	clipboard = SDL_DisplayFormat (image);
	SDL_FreeSurface (image);
}
//---------------------------------------------------------
static void clipboardKopieren (Sint16 x, Sint16 y, Uint16 b, Uint16 h)
{
	SDL_Rect src, dest;
	Uint32 rmask, gmask, bmask, amask;

	if (clipboard != NULL)
		SDL_FreeSurface (clipboard);
	#if SDL_BYTEORDER == SDL_BIG_ENDIAN
		rmask = 0xff000000;
		gmask = 0x00ff0000;
		bmask = 0x0000ff00;
		amask = 0x000000ff;
	#else
		rmask = 0x000000ff;
		gmask = 0x0000ff00;
		bmask = 0x00ff0000;
		amask = 0xff000000;
	#endif
	clipboard = SDL_CreateRGBSurface(SDL_HWSURFACE, (int) b, (int) h, 32, rmask, gmask, bmask, amask);
	if (clipboard == NULL) {
		// printf("Neues Clipboard konnte nicht erzeugt werden!\n");
		return;
	}
	src.x = x;
	src.y = y;
	src.w = b;
	src.h = h;
	dest.x = 0;
	dest.y = 0;
	dest.w = b;
	dest.h = h;
	SDL_BlitSurface(screen, &src, clipboard, &dest);
	SDL_UpdateRect (clipboard, 0, 0, 0, 0);
}
//---------------------------------------------------------
static void clipboardEinfuegen (Sint16 x, Sint16 y)
{
	SDL_Rect src, dest;
	src.x = 0;
	src.y = 0;
	src.w = clipboard->w;
	src.h = clipboard->h;
	dest.x = x;
	dest.y = y;
	dest.w = clipboard->w;
	dest.h = clipboard->h;
	SDL_BlitSurface(clipboard, &src, screen, &dest);
	if (updateOn)
		SDL_UpdateRect (screen, x, y, dest.w, dest.h);
}
//---------------------------------------------------------
static void archivieren ()
{
	SDL_Rect src, dest;
	src.x = 0;
	src.y = 0;
	src.w = screen->w;
	src.h = screen->h;
	dest = src;
	SDL_BlitSurface(screen, &src, archiv, &dest);
	SDL_UpdateRect(archiv, 0,0,0,0);
}
//----------------------------------------------------------
static void restaurieren (Sint16 x, Sint16 y, Uint16 b, Uint16 h)
{
	SDL_Rect src, dest;
	src.x = x;
	src.y = y;
	src.w = b;
	src.h = h;
	dest = src;
	SDL_BlitSurface(archiv, &src, screen, &dest);
	if (updateOn)
		SDL_UpdateRect (screen, x, y, b, h);
}
//---------------------------------------------------------------
static Uint32 tastaturLesen1 ()
{
	tasteLesen = 1;
	while (tasteLesen)
	{
		SDL_Delay (5);
	}
	return ((Uint32) taste)*256*256 + ((Uint32) gedrueckt)*256*256*256*128+ ((Uint32) tiefe);
}
//-------------------------------------------------------------
static void tastaturpufferAn () {
	t_pufferkopf = 0;
	t_pufferende = 0;
	tastaturpuffer = 1;
}
//-------------------------------------------------------------
static void tastaturpufferAus () {
	tastaturpuffer = 0;
}
//-------------------------------------------------------------
static Uint32 tastaturpufferLesen1 ()
{
	Uint32 erg;
	while (t_pufferende == t_pufferkopf)
	{
		SDL_Delay (5);
	}
	erg = t_puffer[t_pufferkopf];
	t_pufferkopf++; //Überlauf von 255 auf 0 automatisch, da Uint8
	return erg;
}
//-------------------------------------------------------------
static Uint32 mausLesen1 ()
{
	mausLesen = 1;
	while (mausLesen)
	{
		SDL_Delay (5);
	}
	return ((Uint32) mausTaste)*256*256*256 + (((Uint32) mausX) << 12) + ((Uint32) mausY);
}
//--------------------------------------------------------------
static void mauspufferAn () {
	m_pufferkopf = 0;
	m_pufferende = 0;
	mauspuffer = 1;
}
//-------------------------------------------------------------
static void mauspufferAus () {
	mauspuffer = 0;
}
//-------------------------------------------------------------
static Uint32 mauspufferLesen1 ()
{
	Uint32 erg;
	while (m_pufferende == m_pufferkopf)
	{
		SDL_Delay (5);
	}
	erg = m_puffer[m_pufferkopf];
	m_pufferkopf++; //Überlauf von 255 auf 0 automatisch, da Uint8
	return erg;
}
//-------------------------------------------------------------
*/
import "C"

import (
	"time"
	"unsafe"
)

var grafikschloss = make(chan int, 1)
var tastaturschloss = make(chan int, 1)
var fensterschloss = make(chan int, 1)
var mausschloss = make(chan int, 1)
var fensterbreite, fensterhoehe uint16

// Es gibt 4 Tastenbelegungen: Standard, SHIFT, ALT GR, ALT GR mit SHIFT.
var z1 [4]string = [4]string{",-.", ";_:", "·–…", "×—÷"}
var z2 [4]string = [4]string{"0123456789", "=!\"§$%&/()", "}¹²³¼½¬{[]", "°¡⅛£¤⅜⅝⅞™±"}
var z3 [4]string = [4]string{"abcdefghijklmnopqrstuvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "æ“¢ð€đŋħ→̣ĸłµ”øþ@¶ſŧ↓„ł«»←", "Æ‘©Ð€ªŊĦı˙&Łº’ØÞΩ®ẞŦ↑‚Ł‹›¥"}
var z4 [4]string = [4]string{",/*-+", ",/*-+", ",/*-+", ",/*-+"}                     //NUM-Block
var z5 [4]string = [4]string{"0123456789", "0123456789", "0123456789", "0123456789"} //NUM-BLOCK
var taste_belegung [4][320]rune                                                      //vier Belegungen pro Taste
//-----------------------------------------------------------------------------

func lock() {
	grafikschloss <- 1
}

func unlock() {
	<-grafikschloss
}

func t_lock() {
	tastaturschloss <- 1
}

func t_unlock() {
	<-tastaturschloss
}

func m_lock() {
	mausschloss <- 1
}

func m_unlock() {
	<-mausschloss
}

func Sperren() {
	fensterschloss <- 1
}

func Entsperren() {
	<-fensterschloss
}

func Fenster(breite, hoehe uint16) {
	lock()
	if fensterZu() {
		if breite > 1920 {
			breite = 1920
		}
		if hoehe > 1200 {
			hoehe = 1200
		}
		fensterhoehe = hoehe
		fensterbreite = breite
		go C.GrafikfensterAn(C.Uint16(breite), C.Uint16(hoehe))
		for !FensterOffen() {
			time.Sleep(100 * 1000 * 1000) //Unter Windows notwendig!!
		}
	}
	unlock()
}

func FensterOffen() bool {
	return uint8(C.FensterOffen()) == 1
}

func fensterZu() bool {
	return uint8(C.FensterZu()) == 1
}

func FensterAus() {
	lock()
	if FensterOffen() {
		C.GrafikfensterAus()
		for !fensterZu() {
			time.Sleep(100 * 1000 * 1000)
		}
	}
	unlock()
}

func Punkt(x, y uint16) {
	lock()
	C.zeichnePunkt(C.Sint16(x), C.Sint16(y))
	unlock()
}

func GibPunktfarbe(x, y uint16) (r, g, b uint8) {
	lock()
	pixel := uint32(C.gibPixel(C.Sint16(x), C.Sint16(y)))
	r = uint8(pixel >> 16)
	g = uint8(pixel >> 8)
	b = uint8(pixel)
	unlock()
	return
}

func Kreis(x, y, r uint16) {
	lock()
	C.zeichneKreis(C.Sint16(x), C.Sint16(y), C.Sint16(r), 0)
	unlock()
}

func Vollkreis(x, y, r uint16) {
	lock()
	C.zeichneKreis(C.Sint16(x), C.Sint16(y), C.Sint16(r), 1)
	unlock()
}

func Ellipse(x, y, rx, ry uint16) {
	lock()
	C.zeichneEllipse(C.Sint16(x), C.Sint16(y), C.Sint16(rx), C.Sint16(ry), 0)
	unlock()
}

func Vollellipse(x, y, rx, ry uint16) {
	lock()
	C.zeichneEllipse(C.Sint16(x), C.Sint16(y), C.Sint16(rx), C.Sint16(ry), 1)
	unlock()
}

func Cls() {
	lock()
	C.clearscreen()
	unlock()
}

func Stiftfarbe(r, g, b uint8) {
	lock()
	C.stiftfarbe(C.Uint8(r), C.Uint8(g), C.Uint8(b))
	unlock()
}

func Linie(x1, y1, x2, y2 uint16) {
	lock()
	C.zeichneStrecke(C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2))
	unlock()
}

func Vollkreissektor(x, y, r, w1, w2 uint16) {
	lock()
	C.kreissektor(C.Sint16(x), C.Sint16(y), C.Sint16(r), 360-C.Sint16(w2), 360-C.Sint16(w1), 1)
	unlock()
}

func Kreissektor(x, y, r, w1, w2 uint16) {
	lock()
	C.kreissektor(C.Sint16(x), C.Sint16(y), C.Sint16(r), 360-C.Sint16(w2), 360-C.Sint16(w1), 0)
	unlock()
}

func Rechteck(x1, y1, b, h uint16) {
	lock()
	C.rechteck(C.Sint16(x1), C.Sint16(y1), C.Sint16(b), C.Sint16(h), 0)
	unlock()
}

func Vollrechteck(x1, y1, b, h uint16) {
	lock()
	C.rechteck(C.Sint16(x1), C.Sint16(y1), C.Sint16(b), C.Sint16(h), 1)
	unlock()
}

func Dreieck(x1, y1, x2, y2, x3, y3 uint16) {
	lock()
	C.dreieck(C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Sint16(x3), C.Sint16(y3), 0)
	unlock()
}

func Volldreieck(x1, y1, x2, y2, x3, y3 uint16) {
	lock()
	C.dreieck(C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Sint16(x3), C.Sint16(y3), 1)
	unlock()
}

func LadeBild(x, y uint16, s string) {
	lock()
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.ladeBild(C.Sint16(x), C.Sint16(y), cs)
	unlock()
}

func SpieleSound(s string) {
	lock()
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	erg := int(C.LoadAndPlaySound(cs))
	if erg == 1 {
		println("Es werden schon die max. Anzahl an Sounds abgespielt!")
	}
	if erg == 2 {
		println("Konnte Sounddatei nicht laden! --> ", s)
	}
	unlock()
}

// INTERN
// Vor.: data stellt die Bytefolge einer WAVE-Datei dar.
//       wartezeit ist die Abwartezeit nach dem Anspielen der WAV-Datei in ms.
// Eff.: Die 'WAV-Datei' wird bzw. ist gerade abgespielt. Der Programmablauf
//       ist dafür um wartezeit ms verzögert worden.
func spieleRAMWAV(data []byte, wartezeit uint32) {
	lock()
	erg := int(C.LadeUndSpieleNote(unsafe.Pointer(&data[0]), C.int(len(data))))
	if erg == 1 {
		println("Es werden schon die max. Anzahl an Sounds abgespielt!")
	}
	if erg == 2 {
		println("Die Daten entsprechen keiner WAV-Datei! Daten nicht geladen!")
	}
	unlock()
	time.Sleep(time.Duration(int64(wartezeit) * 1e6))
}

func LadeBildInsClipboard(s string) {
	lock()
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.ladeBildInsClipboard(cs)
	unlock()
}

func Clipboard_kopieren(x, y, b, h uint16) {
	lock()
	C.clipboardKopieren(C.Sint16(x), C.Sint16(y), C.Uint16(b), C.Uint16(h))
	unlock()
}

func Clipboard_einfuegen(x, y uint16) {
	lock()
	C.clipboardEinfuegen(C.Sint16(x), C.Sint16(y))
	unlock()
}

func Archivieren() {
	lock()
	C.archivieren()
	unlock()
}

func Restaurieren(x1, y1, b, h uint16) {
	lock()
	C.restaurieren(C.Sint16(x1), C.Sint16(y1), C.Uint16(b), C.Uint16(h))
	unlock()
}

func Schreibe(x, y uint16, s string) {
	lock()
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.schreibe(C.Sint16(x), C.Sint16(y), cs)
	unlock()
}

func SetzeFont(s string, groesse int) (erg bool) {
	lock()
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	if int(C.setFont(cs, C.int(groesse))) == 0 {
		erg = true
	} else {
		erg = false
	}
	unlock()
	return
}

func GibFont() (erg string) {
	lock()
	cs := C.getFont()
	//defer C.free(unsafe.Pointer(cs))
	erg = C.GoString(cs)
	unlock()
	return
}

func SchreibeFont(x, y uint16, s string) {
	lock()
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	if int(C.write(C.Sint16(x), C.Sint16(y), cs)) == 1 {
		println("FEHLER: Kein aktueller Font: ", C.GoString(C.getFont()))
	}
	unlock()
}

func UpdateAus() {
	lock()
	C.updateAus()
	unlock()
}

func UpdateAn() {
	lock()
	C.updateAn()
	unlock()
}

func TastaturLesen1() (taste uint16, gedrueckt uint8, tiefe uint16) {
	var tastenwert uint32
	t_lock()
	tastenwert = uint32(C.tastaturLesen1())
	t_unlock()
	tiefe = uint16(tastenwert % 65536)
	tastenwert = tastenwert >> 16
	gedrueckt = uint8(tastenwert >> 15)
	taste = uint16(tastenwert % 32768) //oberstes Bit rausschieben
	return
}

func TastaturpufferAn() {
	t_lock()
	C.tastaturpufferAn()
	t_unlock()
}

func TastaturpufferAus() {
	t_lock()
	C.tastaturpufferAus()
	t_unlock()
}

func TastaturpufferLesen1() (taste uint16, gedrueckt uint8, tiefe uint16) {
	var tastenwert uint32
	t_lock()
	tastenwert = uint32(C.tastaturpufferLesen1())
	t_unlock()
	tiefe = uint16(tastenwert % 65536)
	tastenwert = tastenwert >> 16
	gedrueckt = uint8(tastenwert >> 15)
	taste = uint16(tastenwert % 32768)
	return
}

func Tastaturzeichen(taste, tiefe uint16) rune {
	switch tiefe {
	case 0, 4096, 8192 + 1, 8192 + 2, 8192 + 3, 4096 + 8192 + 1, 4096 + 8192 + 2, 4096 + 8192 + 3: // kein SHIFT, kein ALT GR, NUMLOCK an oder aus, CAPSLOCK an mit SHIFT
		return taste_belegung[0][taste]
	case 1, 2, 3, 4096 + 1, 4096 + 2, 4096 + 3, 8192, 4096 + 8192: // SHIFT, kein ALT GR, NUMLOCK an oder aus, CAPSLOCK an ohne SHIFT
		return taste_belegung[1][taste]
	case 16384, 16384 + 4096, 16384 + 8192 + 1, 16384 + 8192 + 2, 16384 + 8192 + 3, 16384 + 8192 + 4096 + 1, 16384 + 8192 + 4096 + 2, 16384 + 8192 + 4096 + 3: // kein SHIFT, ALT GR, NUMLOCK an oder aus, CAPSLOCK an mit SHIFT
		return taste_belegung[2][taste]
	case 16384 + 1, 16384 + 2, 16384 + 3, 16384 + 4096 + 1, 16384 + 4096 + 2, 16384 + 4096 + 3, 16384 + 8192, 16384 + 8192 + 4096: // ALT GR und SHIFT, NUMLOCK an oder aus, CAPSLOCK an ohne SHIFT
		return taste_belegung[3][taste]
	default:
		return 0
	}
}

func MausLesen1() (taste uint8, status int8, mausX, mausY uint16) {
	var tastenwert uint32
	m_lock()
	tastenwert = uint32(C.mausLesen1())
	m_unlock()
	taste = uint8(tastenwert >> 24)
	if taste < 64 {
		status = 0 //Zustand wird gehalten
	} else if taste > 128 {
		status = 1 //gerade gedrückt
		taste = taste - 128
	} else { //zwischen 64 und 128
		status = -1 //gerade losgelassen
		taste = taste - 64
	}
	mausY = uint16(tastenwert % 4096)
	tastenwert = tastenwert >> 12
	mausX = uint16(tastenwert % 4096)
	return
}

func MauspufferAn() {
	m_lock()
	C.mauspufferAn()
	m_unlock()
}

func MauspufferAus() {
	m_lock()
	C.mauspufferAus()
	m_unlock()
}

func MauspufferLesen1() (taste uint8, status int8, mausX, mausY uint16) {
	var tastenwert uint32
	m_lock()
	tastenwert = uint32(C.mauspufferLesen1())
	m_unlock()
	taste = uint8(tastenwert >> 24)
	if taste < 64 {
		status = 0 //Zustand wird gehalten
	} else if taste > 128 {
		status = 1 //gerade gedrückt
		taste = taste - 128
	} else { //zwischen 64 und 128
		status = -1 //gerade losgelassen
		taste = taste - 64
	}
	mausY = uint16(tastenwert % 4096)
	tastenwert = tastenwert >> 12
	mausX = uint16(tastenwert % 4096)
	return
}

func Grafikzeilen() uint16 {
	return fensterhoehe
}

func Grafikspalten() uint16 {
	return fensterbreite
}

func init() {
	// Es folgt die Initalisierung der Tastaturbelegung auf Deutsch.
	// Das wird für die Funktion 'Tastaturzeichen(taste, tiefe) rune' benötigt.
	for i := 0; i < 4; i++ {
		index := 0
		for _, e := range z1[i] {
			taste_belegung[i][index+44] = e
			index++
		}
		index = 0
		for _, e := range z2[i] {
			taste_belegung[i][index+48] = e
			index++
		}
		index = 0
		for _, e := range z5[i] {
			taste_belegung[i][index+256] = e
			index++
		} //Num-Block
		index = 0
		for _, e := range z4[i] {
			taste_belegung[i][index+266] = e
			index++
		} //Num-Block
		index = 0
		for _, e := range z3[i] {
			taste_belegung[i][index+97] = e
			index++
		}
	}
	// kein SHIFT, kein ALT GR
	taste_belegung[0][43] = '+'
	taste_belegung[0][35] = '#'
	taste_belegung[0][252] = 'ü'
	taste_belegung[0][246] = 'ö'
	taste_belegung[0][228] = 'ä'
	taste_belegung[0][223] = 'ß'
	taste_belegung[0][180] = '´'
	taste_belegung[0][94] = '^'
	taste_belegung[0][60] = '<'
	taste_belegung[0][32] = ' '
	// SHIFT, kein ALT GR
	taste_belegung[1][43] = '*'
	taste_belegung[1][35] = '\''
	taste_belegung[1][252] = 'Ü'
	taste_belegung[1][246] = 'Ö'
	taste_belegung[1][228] = 'Ä'
	taste_belegung[1][223] = '?'
	taste_belegung[1][180] = '`'
	taste_belegung[1][94] = '°'
	taste_belegung[1][60] = '>'
	taste_belegung[1][32] = ' '
	// kein SHIFT, ALT GR
	taste_belegung[2][43] = '~'
	taste_belegung[2][35] = '`'
	taste_belegung[2][252] = '¨'
	taste_belegung[2][246] = '˝'
	taste_belegung[2][228] = '^'
	taste_belegung[2][223] = '\\'
	taste_belegung[2][180] = '¸'
	taste_belegung[2][94] = '¬'
	taste_belegung[2][60] = '|'
	taste_belegung[2][32] = ' '
	// SHIFT, ALT GR
	taste_belegung[3][43] = '¯'
	taste_belegung[3][35] = '`'
	taste_belegung[3][252] = '¨'
	taste_belegung[3][246] = '˝'
	taste_belegung[3][228] = '^'
	taste_belegung[3][223] = '¿'
	taste_belegung[3][180] = '¸'
	taste_belegung[3][94] = '¬'
	taste_belegung[3][60] = '¦'
	taste_belegung[3][32] = ' '
}

//------------------------ENDE----------------------------------
