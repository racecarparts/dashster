function openUrl(url) {
    // bound function in service/bindings.go
    open(url)
}

function timeIntervalStr(intervalMillis) {
    let d = Date.now()
    let now = new Date(d)
    let next = new Date(d + intervalMillis)
    let timeOpts = {
        hour: 'numeric',
        minute: 'numeric'
    }
    let tf = new Intl.DateTimeFormat([], timeOpts)
    return tf.format(now) + " -> " + tf.format(next)
}

function writeWidget(dataContainerName, innerHTML) {
    const continerSelector = "[data-container=" + dataContainerName + "]"
    const containerEl = document.querySelectorAll(continerSelector)
    if (containerEl.length !== 1) {
        console.log("Invalid container to write data: " + continerSelector)
        return
    }
    let container = containerEl.item(0);
    container.innerHTML = innerHTML
}

$( document ).ready(function() {
    var heights = $(".panel").map(function() {
            return $(this).height();
        }).get(),

        maxHeight = Math.max.apply(null, heights);

    $(".panel").height(maxHeight);
});