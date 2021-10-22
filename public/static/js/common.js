function openUrl(url) {
    // bound function in service/bindings.go
    open(url)
}

// async function fetchWithTimeout(resource, options = {}) {
//     const { timeout = 240000 } = options;
//
//     const controller = new AbortController();
//     const id = setTimeout(() => controller.abort(), timeout);
//     const response = await fetch(resource, {
//         ...options,
//         signal: controller.signal
//     });
//     clearTimeout(id);
//     return response;
// }

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

function getButtonEl(dataButtonName) {
    const buttonSelector = "[data-button-name=" + dataButtonName + "]"
    const buttonEl = document.querySelectorAll(buttonSelector)
    if (buttonEl.length !== 1) {
        console.log("Invalid button: " + buttonSelector)
        return
    }
    return buttonEl.item(0);
}

function disableButton(dataButtonName) {
    getButtonEl(dataButtonName).disabled = true
}

function enableButton(dataButtonName) {
    getButtonEl(dataButtonName).disabled = false
}

function writeWidget(dataContainerName, innerHTML) {
    const containerSelector = "[data-container=" + dataContainerName + "]"
    const containerEl = document.querySelectorAll(containerSelector)
    if (containerEl.length !== 1) {
        console.log("Invalid container to write data: " + containerSelector)
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