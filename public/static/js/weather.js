const weatherInterval = 3600000 // every hour

function writeWeather(weatherData) {
    let weather = ""
    for (let i = 0; i < weatherData.length; i++) {
        let location = weatherData[i].location.padEnd(16) + "  "
        let cond = weatherData[i].cur_conditions.padEnd(16) + "  "
        let temp = weatherData[i].cur_temp.padEnd(5) + "\n"
        weather += location + cond + temp
    }

    writeWidget('weather-interval', timeIntervalStr(weatherInterval))
    writeWidget('weather', weather)
}

function loadWeather() {
    writeWidget('weather-interval', 'loading...')
    fetch('/weather', {
        method: 'get'
    })
        .then(r => r.json())
        .then(jsonData => {
            writeWeather(jsonData)
        })
        .catch(err => {
            writeWidget('weather-interval', 'error: ' + err)
            console.log(err)
        })
}

loadWeather()
setInterval(() => {
    loadWeather()
}, weatherInterval)