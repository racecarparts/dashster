const agendaInterval = 900000  // every 15 minutes

function writeMyCal(myCalData) {
    writeWidget("mycal-interval", timeIntervalStr(agendaInterval))
    writeWidget('my-cal', myCalData.agenda)
}

function loadMyCal() {
    writeWidget('mycal-interval', 'loading...')
    fetch('/mycal', {
        method: 'get'
    })
        .then(r => r.json())
        .then(jsonData => {
            writeMyCal(jsonData)
        })
        .catch(err => {
            writeWidget('pr-interval', 'error: ' + err)
            console.log(err)
        })
}

loadMyCal()
setInterval(() => {
    loadMyCal()
}, agendaInterval)