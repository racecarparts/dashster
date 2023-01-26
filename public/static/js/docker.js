const dockerInterval = 15000 // every 15 seconds
let dockerIntervalId

function writeDocker(dockerData) {
    writeWidget('docker', dockerData.stat)
}

function loadDocker() {
    fetch('/docker', {
        method: 'get'
    })
        .then(r => r.json())
        .then(jsonData => {
            writeDocker(jsonData)
            dockerIntervalId = setupInterval(dockerIntervalId, dockerInterval, loadDocker)
        })
        .catch(err => {
            console.log(err)
        })
}

loadDocker()
