const dockerInterval = 30000 // every 30 seconds

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
        })
        .catch(err => {
            console.log(err)
        })
}

loadDocker()
setInterval(() => {
    loadDocker()
}, dockerInterval)