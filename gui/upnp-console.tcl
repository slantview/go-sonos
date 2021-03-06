#!/usr/bin/wish

package require Tk
package require json
package require http

namespace eval state {
	variable volume 0 mute 0
}

set URL "http://localhost:8080"
set CONTROL "$URL/control"
set BROWSE "$URL/browse"

proc on_set_volume {token} {
	upvar $token state
	puts $state(code)
}

proc on_volume_scale {value} {
	http::geturl "$::CONTROL?method=set-volume&value=$state::volume" -command on_set_volume
}

proc on_get_volume {token} {
	upvar $token state
	set volume [dict get [json::json2dict $state(body)] Value]
	puts "got volume $volume"
	set state::volume $volume
}

proc get_volume {} {
	http::geturl "$::CONTROL?method=get-volume" -command on_get_volume
}

proc on_mute_button {} {
	puts "set mute state $state::mute"
}

proc on_get_genres {token} {
	upvar $token state
	set msg [dict get [json::json2dict $state(body)] Value ]
	foreach item $msg {
		set title [dict get $item Title]
		lappend genres $title
	}
	.top_level.notebook.content.genre configure -values $genres
}

proc get_genres {} {
	http::geturl "$::BROWSE?method=get-all-genres" -command on_get_genres
}

wm title . "UPnP Console"
wm minsize . 600 600
pack [frame .top_level] -fill both
pack [ttk::notebook .top_level.notebook] -fill both
pack [frame .top_level.notebook.control] -fill both
pack [frame .top_level.notebook.content] -fill both
.top_level.notebook add .top_level.notebook.control -text "Control"
.top_level.notebook add .top_level.notebook.content -text "Content"

set vf .top_level.notebook.control.volume
pack [labelframe $vf -text "Volume"] -fill both
pack [label $vf.volume_label -text "Volume"]
pack [scale $vf.volume_scale -command on_volume_scale -from 0 -to 100 -variable state::volume -showvalue 1 -orient horizontal]
pack [label $vf.mute_label -text "Mute"]
pack [checkbutton $vf.mute -command on_mute_button -variable state::mute]

set vf .top_level.notebook.content
pack [ttk::combobox $vf.genre -text "Genre"]
pack [ttk::combobox $vf.artist -text "Artist"]
pack [ttk::combobox $vf.album -text "Album"]
pack [ttk::combobox $vf.track -text "Trackk"]

get_volume
get_genres

