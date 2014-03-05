import QtQuick 2.0
import Ubuntu.Components 0.1
import Ubuntu.Components.ListItems 0.1

MainView {
    applicationName: "JenkinsMonitor"
    width: 640
    height: 480

    ListView {
        id: serverView
        width: 120;
        model: servers.len
        delegate: serverDelegate
        anchors.fill: parent
    }

    Component {
        id: serverDelegate
        Column {
            property var currentServer: servers.server(index)
            Text {
                text: "Server: " + currentServer.address
            }
            Text {
                text: "Port: " + currentServer.port
            }
            // Header {
            //     text: "Jobs"
            // }
            Repeater {
                anchors.top: parent.bottom
                model: currentServer.jobs.len
                delegate: JobComponent {
                    job: currentServer.jobs.job(index)
                }
            }
        }
    }
}
