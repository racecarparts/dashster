package service

import (
    "github.com/webview/webview"
)

func RegisterBindings(w webview.WebView) webview.WebView {
    w.Bind("open", openUrlInNewWindow)
    //w.Bind("worldClock", worldClock)
    //w.Bind("Calendar", Calendar)
    //w.Bind("weather", weather)
    //w.Bind("docker", docker)
    //w.Bind("pullRequests", pullRequests)

    return w
}
