import QtQuick 2.0
import Ubuntu.Components 0.1
import Ubuntu.Components.ListItems 0.1 as ListItems

Text {
	property var job

	width: 400
	text: job.name
	Rectangle {
		color: parent.job.renderColor()
		anchors.fill: parent
	}
}