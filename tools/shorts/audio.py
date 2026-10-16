"""Local Chinese narration and original synthesized music; no external service."""
import json,subprocess,sys,wave
from pathlib import Path
import numpy as np

root=Path(__file__).resolve().parent
out=Path(sys.argv[1]).resolve();out.mkdir(parents=True,exist_ok=True)
story=json.loads((root/'story.json').read_text())
sr=48000;duration=story['duration'];voice=np.zeros(int(duration*sr));music=np.zeros_like(voice)

def command(*args):subprocess.run(args,check=True,stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL)
def stamp(seconds):
    ms=round(seconds*1000);return f'{ms//3600000:02}:{ms//60000%60:02}:{ms//1000%60:02},{ms%1000:03}'
subtitles=[];mixes=[]
for i,row in enumerate(story['segments']):
    text=out/f'voice-{i}.txt';aiff=out/f'voice-{i}.aiff';wav=out/f'voice-{i}.wav'
    text.write_text(row['voice'])
    command('say','-v','Tingting','-r','210','-o',str(aiff),'-f',str(text))
    seconds=float(subprocess.check_output(['ffprobe','-v','quiet','-show_entries','format=duration','-of','csv=p=0',str(aiff)],text=True))
    slot=row['end']-row['start']-.18
    tempo=max(1,seconds/slot)
    if tempo>1.45:raise RuntimeError(f'Segment {i} narration is too dense: tempo={tempo:.2f}')
    command('ffmpeg','-y','-v','error','-i',str(aiff),'-af',f'atempo={tempo:.5f},highpass=f=80,loudnorm=I=-18:TP=-2:LRA=7','-ar',str(sr),'-ac','1','-c:a','pcm_s16le',str(wav))
    with wave.open(str(wav),'rb') as f:x=np.frombuffer(f.readframes(f.getnframes()),dtype='<i2').astype(float)/32768
    start=round((row['start']+.04)*sr);x=x[:len(voice)-start];voice[start:start+len(x)]+=x
    mixes.append({'segment':i,'speech_seconds':seconds,'tempo':round(tempo,4)})
    subtitles.append(f'{i+1}\n{stamp(row["start"])} --> {stamp(row["end"])}\n{row["caption"]}\n')
# Original 120 BPM pulse: quiet enough for a phone speaker under narration.
for beat,t0 in enumerate(np.arange(0,duration,.5)):
    n=int(.36*sr);t=np.arange(n)/sr
    tone=[261.63,329.63,392.0,293.66][(beat//8)%4]*([1,2,1.5,2][beat%4])
    x=.035*np.sin(2*np.pi*tone*t)*np.exp(-t*13)
    if beat%2==0:x+=.06*np.sin(2*np.pi*(55*t+35*(1-np.exp(-t*24))/24))*np.exp(-t*22)
    start=round(t0*sr);n=min(n,len(music)-start);music[start:start+n]+=x[:n]
for row in story['segments']:
    start=int(row['start']*sr);n=min(int(.13*sr),len(music)-start);t=np.arange(n)/sr
    music[start:start+n]+=.04*np.sin(2*np.pi*(900*t-1500*t*t))*np.exp(-t*40)
mix=voice+music
mix*=np.minimum(1,np.arange(len(mix))/(sr*.03))*np.minimum(1,(len(mix)-np.arange(len(mix)))/(sr*.24))
peak=np.max(np.abs(mix));mix=mix/min(1e6,max(1,peak/.88))
with wave.open(str(out/'mix.wav'),'wb') as f:f.setnchannels(1);f.setsampwidth(2);f.setframerate(sr);f.writeframes((mix*32767).astype('<i2').tobytes())
(out/'PasteFerry-2.1-douyin.srt').write_text('\n'.join(subtitles))
(out/'narration-timing.json').write_text(json.dumps(mixes,indent=2)+'\n')
print('Narration, original soundtrack and captions generated')
