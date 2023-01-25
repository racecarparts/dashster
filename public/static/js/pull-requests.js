const prInterval = 900000 // every 15 minutes
let prIntervalId;

function writePRRows(prDataArr, outerEl, sectionTitle) {
    for (let i = 0; i < prDataArr.length; i++) {
        let pr = prDataArr[i]
        let rowEl = document.createElement("tr")

        let repoNameHeadEl = document.createElement("th")
        repoNameHeadEl.setAttribute("scope", "row")
        repoNameHeadEl.className = "text-nowrap"
        repoNameHeadEl.innerHTML = pr.repository_name
        rowEl.appendChild(repoNameHeadEl)

        // let prShaEl = document.createElement("th")
        // prShaEl.className = "text-nowrap"
        // prShaEl.innerHTML = pr.sha
        // rowEl.appendChild(prShaEl)

        let prLinkCellEl = document.createElement("td")
        let prLinkEl = document.createElement("a")
        prLinkEl.setAttribute("onclick", "openUrl('"+ pr.web_url +"')")
        prLinkEl.setAttribute("href", "#")
        prLinkEl.innerHTML = pr.number
        prLinkCellEl.appendChild(prLinkEl)
        rowEl.appendChild(prLinkCellEl)

        let prUserEl = document.createElement("th")
        prUserEl.className = "text-nowrap"
        prUserEl.innerHTML = pr.user
        rowEl.appendChild(prUserEl)

        let titleEl = document.createElement("td")
        titleEl.innerHTML = pr.title
        rowEl.appendChild(titleEl)

        let approvalsEl = document.createElement("td")
        approvalsEl.className = "text-nowrap"
        let approvals = ""
        for (let j = 0; j < pr.reviews.length; j++) {
            let review = pr.reviews[j]
            approvals += review.user.login + ": " + review.state + "<br>"
        }
        approvalsEl.innerHTML = approvals
        rowEl.appendChild(approvalsEl)

        let updatedEl = document.createElement("td")
        updatedEl.className = "text-nowrap"
        updatedEl.innerHTML = pr.updated_at
        rowEl.appendChild(updatedEl)

        outerEl.appendChild(rowEl)
    }
    return outerEl
}

function writePRs(prData) {
    if (!prData.my_prs && !prData.requested_prs) {
        writeWidget('pr-interval', timeIntervalStr(prInterval))
        return
    }
    let outerEl = document.createElement("div")
    if (prData.message.length > 0) {
        writeWidget('pull-requests', prData.message)
    }
    outerEl = writePRRows(prData.my_prs, outerEl, "My PRs")
    outerEl = writePRRows(prData.requested_prs, outerEl, "Requested PRs")

    writeWidget('pr-interval', timeIntervalStr(prInterval))
    writeWidget('pull-requests', outerEl.innerHTML.toString())
}

function loadPRs() {
    disableButton('load-prs')
    writeWidget('pr-interval', 'loading...')
    fetch('/mergerequests', {
        method: 'get'
    })
        .then(r => r.json())
        .then(jsonData => {
            writePRs(jsonData)
            enableButton('load-prs')
            
            prIntervalId = setupInterval(prIntervalId, prInterval, loadPRs)
        })
        .catch(err => {
            writeWidget('pr-interval', 'error: ' + err)
            enableButton('load-prs')
            console.log(err)
        })
}

// writePRs(
//     {message: "",
//         my_prs: [
//             {
//                 is_draft: false,
//                 number: 0,
//                 user: "racecarparts",
//                 title: "something cool",
//                 repository_name: "whirled-peas",
//                 review_url: "https://github.com/racecarparts/whirled-peas/pull/0",
//                 reviews: [{user: {login: "racecarparts"}, state: "APPROVED"}]
//             },
//         ],
//         requested_prs:[
//             {
//                 is_draft: false,
//                 number: 56,
//                 user: "sviatkh",
//                 title: "TE-28583-setup-celery-on-dev-server",
//                 repository_name: "teem-dev-deploy",
//                 review_url: "https://github.com/enderlabs/teem-dev-deploy/pull/56",
//                 reviews: [{user: {login: "racecarparts"}, state: "APPROVED"}],
//                 sha: "6dcb09b"
//             },
//             {
//                 is_draft: false,
//                 number: 276,
//                 user: "jackson-david",
//                 title: "TE-23320 Subscription Renewals",
//                 repository_name: "reservations",
//                 review_url: "https://github.com/enderlabs/reservations/pull/276",
//                 reviews: [
//                     {user: {login: "jacobdh"}, state: "APPROVED"},
//                     {user: {login: "jackson-david"}, state: "COMMENTED"},
//                     {user: {login: "jackson-david"}, state: "COMMENTED"}
//                     ]
//             },
//             {
//                 is_draft: false,
//                 number: 6356,
//                 user: "jacobdh",
//                 title: "TE-29630: Releasefix - Fix diagnostics redis cluster mode error",
//                 repository_name: "eventboard.io",
//                 review_url: "https://github.com/enderlabs/eventboard.io/pull/6356",
//                 reviews: []
//             }
//         ]
//     })

loadPRs()
