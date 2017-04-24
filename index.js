'use strict'

var wavesurfer = WaveSurfer.create({
    container: '#waveform',
    waveColor: 'blue',
    progressColor: 'black'
});

wavesurfer.load("track1.mp3");

var playBtn = document.getElementById("play-track");
playBtn.addEventListener("click", function(ev) {
    wavesurfer.playPause();
}, false);