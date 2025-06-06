import time
from typing import Mapping
from werkzeug import Request, Response
from aiexec_plugin import Endpoint

html = """<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        html, body {
            margin: 0;
            padding: 0;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: white;
            color: black;
            overflow: hidden;
        }

        #svg-container {
            height: 100vh;
            width: 100vw;
        }

        #svg-container svg {
            height: 100%;
            width: 100%;
        }

        #svg-container svg path {
            fill: black;
            stroke: none;
        }

        .overlay {
            position: fixed;
            top: 0;
            left: 0;
            bottom: 0;
            right: 0;
            font-weight: bold;
            text-align: center;
        }

        .overlay .controls {
            margin: 10px;
            background-color: #ffffff75;
            border-radius: 10px;
            width: min-content;
        }

        .overlay .controls .control {
            padding: 5px;
        }

        #loading-text {
            position: relative;
            top: 50%;
            font-size: large;
        }
    </style>
    <title>Bad Apple!!.SVG</title>
</head>
<body>
    <div class="overlay">
        <span id="loading-text">Loading...</span>
        <div class="controls">
            <div class="control">
                <label for="">Color 1</label>
                <input type="color" name="color1" id="color1" value="#000000">
            </div>
            <div class="control">
                <label for="">Color 2</label>
                <input type="color" name="color2" id="color2" value="#ffffff">
            </div>
        </div>
    </div>
    <div id="svg-container">
        <svg xmlns="http://www.w3.org/2000/svg" version="1.1"
             viewBox="0 0 480 360" width="480" height="360"
             preserveAspectRatio="xMidYMid meet">
        </svg>
    </div>
    <audio id="song" src="https://buzzbyte.github.io/BadApple.SVG/bad-apple.wav" preload="auto"></audio>
    <script>
        window.addEventListener("load", function() {
            const audioContext = new AudioContext();
            const audio = document.getElementById('song');
            const loadingText = document.getElementById("loading-text");

            var svgSupport = !!document.createElementNS && !!document.createElementNS('http://www.w3.org/2000/svg', 'svg').createSVGRect;

            if (!svgSupport) {
                loadingText.innerHTML = "Your browser doesn't support SVG.";
            } else {
                let svgElem = document.querySelector("#svg-container svg");
                let pathElem = document.createElementNS('http://www.w3.org/2000/svg', 'path');
                pathElem.setAttribute('transform', "translate(0.000000,360.000000) scale(0.100000,-0.100000)")
                svgElem.appendChild(pathElem);

                setupControls(pathElem);

                getPathData().then(pathData => {
                    let autoplayAttempt = audio.play();
                    if (autoplayAttempt) {
                        autoplayAttempt.then(() => {
                            playAnim(pathData, pathElem);
                        }).catch(err => {
                            // autoplay probably failed...
                            manualPlay(pathData, pathElem);
                        });
                    } else {
                        // browser's doing dumb shit
                        manualPlay(pathData, pathElem);
                    }
                });
            }
        });

        // Sets up the color controls
        function setupControls(pathElem) {
            let color1 = document.getElementById('color1');
            let color2 = document.getElementById('color2');

            color1.addEventListener("input", function(event) {
                if (pathElem) {
                    pathElem.style.fill = event.target.value;
                }
            });

            color2.addEventListener("input", function(event) {
                let invertedColor = 0xFFFFFF ^ parseInt(event.target.value.substring(1), 16)
                invertedColor = ("000000" + invertedColor.toString(16)).slice(-6);

                document.body.style.backgroundColor = event.target.value;
                document.body.style.color = "#"+ invertedColor;
            });

            color1.select();
            color2.select();
        }

        // For when autoplay is blocked
        function manualPlay(pathData, pathElem) {
            const loadingText = document.getElementById("loading-text"),
                  audio       = document.getElementById("song"),
                  overlay     = document.getElementsByClassName("overlay")[0];
            
            overlay.addEventListener("click", function playClick(event) {
                overlay.removeEventListener("click", playClick);
                loadingText.innerHTML = "";
                overlay.style.cursor = "auto";
                audio.play()
                playAnim(pathData, pathElem);
            });
            overlay.style.cursor = "pointer";
            loadingText.innerHTML = "Click to Play";
        }

        // The main animation controller
        function playAnim(pathData, pathElem) {
            const audio = document.getElementById('song');
            let fps = 15,
                fpsInterval = 1000 / fps,
                frameInterval = (audio.duration / pathData.length) * 1000,
                then = audio.currentTime * 1000,
                start = then,
                frameIndex = 0;
            
            // The animation loop
            function loopAnim() {
                requestAnimationFrame(loopAnim);
                now = audio.currentTime * 1000;
                elapsed = now - then;

                audioSynched = Math.round(now / fpsInterval) + 1;

                if (frameIndex > pathData.length) {
                    frameIndex = pathData.length - 1;
                    cancelAnimationFrame(loopAnim);
                }
                
                if (elapsed > fpsInterval) {
                    then = now - (elapsed % fpsInterval);
                    frameData = pathData[frameIndex];
                    pathElem.setAttribute('d', frameData);

                    // Sync playback with audio if frameIndex is at least
                    // 5 frames behind
                    if (audioSynched > frameIndex + 5) {
                        console.log("Syncing with audio")
                        frameIndex = audioSynched;
                    } else {
                        console.log("Playing smoothly")
                        frameIndex += 1;
                    }
                }
            }

            loopAnim();
        }

        // Gets the list of path data frames from the JSON file
        async function getPathData() {
            let {e, xhttp} = await getRequest("pathdata.json");
            return JSON.parse(xhttp.responseText);
        }

        // Sends a simple get request to url and shows progress
        function getRequest(url) {
            return new Promise((resolve, reject) => {
                var loadingText = document.getElementById("loading-text");
                var xhttp = new XMLHttpRequest();
                xhttp.addEventListener("load", function(event) {
                    loadingText.innerHTML = "";
                    resolve({event, xhttp});
                });
                xhttp.addEventListener("error", function(event) {
                    reject(event);
                });
                xhttp.addEventListener("progress", function(event) {
                    let totalLen = event.lengthComputable ? event.total : 1024 * 1024 * 50;
                    let percent = 100 * (event.loaded / totalLen);
                    if (!event.lengthComputable && percent > 99) {
                        percent = 99;
                    }
                    let percentInt = Math.round(percent);
                    loadingText.innerHTML = `Loading ${percentInt}% ...`;
                });
                xhttp.open("GET", url, true);
                xhttp.send();
            });
        }
    </script>
</body>
</html>"""

class BadappleEndpoint(Endpoint):
    def _invoke(self, r: Request, values: Mapping, settings: Mapping) -> Response:
        """
        Invokes the endpoint with the given request.
        """
        return Response(html, status=200, content_type="text/html")
