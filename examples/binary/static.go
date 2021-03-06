// Code generated by embed. DO NOT EDIT.

package main

import (
	"github.com/inabyte/embed/embedded"
	"unsafe"
)

// FS return file system
var FS embedded.FileSystem

// FileHandler return http file server implements http.Handler
func FileHandler() embedded.Handler {
	return embedded.GetFileServer(FS)
}

var staticData [1366902]byte

func init() {

	bytes := staticData[:]
	str := *(*string)(unsafe.Pointer(&bytes))

	FS = embedded.New(82)

	FS.AddFile( /* /LICENSE.txt */ str[1366848:1366860],
		/* LICENSE.txt */ str[1366849:1366860],
		"",
		17128, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* bkxVFmNed93ddgk0UhjMNeFvChs-gz */ str[1364992:1365022],
		true, bytes[0:5892], str[0:5892])

	FS.AddFile( /* /README.txt */ str[1366871:1366882],
		/* README.txt */ str[1366872:1366882],
		"",
		1128, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* vFLUWdil9cOYtEiGD45xZ8jRrE0-gz */ str[1365142:1365172],
		true, bytes[5892:6565], str[5892:6565])

	FS.AddFile( /* /assets/css/fontawesome-all.min.css */ str[1363193:1363228],
		/* fontawesome-all.min.css */ str[1363205:1363228],
		"",
		55958, 1559863020,
		/* text/css; charset=utf-8 */ str[1366541:1366564],
		/* c4de6La-vqXxiGjpZNoKOx8WJ8Y-gz */ str[1365292:1365322],
		true, bytes[6565:18573], str[6565:18573])

	FS.AddFile( /* /assets/css/images/intro.svg */ str[1366365:1366393],
		/* intro.svg */ str[1366384:1366393],
		"",
		567, 1559863020,
		/* image/svg+xml */ str[1366822:1366835],
		/* t84017DIAYnqQNV1LkzbuQAzS3s-gz */ str[1365442:1365472],
		true, bytes[18573:18844], str[18573:18844])

	FS.AddFile( /* /assets/css/main.css */ str[1366650:1366670],
		/* main.css */ str[1366662:1366670],
		"",
		49030, 1559863020,
		/* text/css; charset=utf-8 */ str[1366541:1366564],
		/* rEn11xiezng9OLc3ecA-vyew5QY-gz */ str[1365592:1365622],
		true, bytes[18844:26565], str[18844:26565])

	FS.AddFile( /* /assets/css/noscript.css */ str[1366517:1366541],
		/* noscript.css */ str[1366529:1366541],
		"",
		572, 1559863020,
		/* text/css; charset=utf-8 */ str[1366541:1366564],
		/* EvGtatGP09xtSKxQZYnxcgICJUs-gz */ str[1365742:1365772],
		true, bytes[26565:26728], str[26565:26728])

	FS.AddFile( /* /assets/js/breakpoints.min.js */ str[1366280:1366309],
		/* breakpoints.min.js */ str[1366291:1366309],
		"",
		2387, 1559863020,
		/* application/javascript */ str[1366564:1366586],
		/* 3xLJWUxs69D9OpPXIhXSqh8Y_bs-gz */ str[1365892:1365922],
		true, bytes[26728:27516], str[26728:27516])

	FS.AddFile( /* /assets/js/browser.min.js */ str[1366444:1366469],
		/* browser.min.js */ str[1366455:1366469],
		"",
		1803, 1559863020,
		/* application/javascript */ str[1366564:1366586],
		/* kK2fIuULYAlNBcSN8UALcB7-9a0-gz */ str[1366042:1366072],
		true, bytes[27516:28315], str[27516:28315])

	FS.AddFile( /* /assets/js/jquery.min.js */ str[1366493:1366517],
		/* jquery.min.js */ str[1366504:1366517],
		"",
		88141, 1559863020,
		/* application/javascript */ str[1366564:1366586],
		/* kIglQh4CAR5PtJFdC6gNOcyx_JQ-gz */ str[1364722:1364752],
		true, bytes[28315:59081], str[28315:59081])

	FS.AddFile( /* /assets/js/jquery.scrollex.min.js */ str[1363869:1363902],
		/* jquery.scrollex.min.js */ str[1363880:1363902],
		"",
		2164, 1559863020,
		/* application/javascript */ str[1366564:1366586],
		/* JXgTd4NdCsvkVAjjgVX_lfiie-w-gz */ str[1366012:1366042],
		true, bytes[59081:59934], str[59081:59934])

	FS.AddFile( /* /assets/js/jquery.scrolly.min.js */ str[1363967:1363999],
		/* jquery.scrolly.min.js */ str[1363978:1363999],
		"",
		770, 1559863020,
		/* application/javascript */ str[1366564:1366586],
		/* -MKQleBRnWpyhKGGnzIbG8obLxQ-gz */ str[1365862:1365892],
		true, bytes[59934:60431], str[59934:60431])

	FS.AddFile( /* /assets/js/main.js */ str[1366670:1366688],
		/* main.js */ str[1366681:1366688],
		"",
		1975, 1559863020,
		/* application/javascript */ str[1366564:1366586],
		/* 7i_emeZbZa9_aX-0C6cWqQ15wdE-gz */ str[1365622:1365652],
		true, bytes[60431:61201], str[60431:61201])

	FS.AddFile( /* /assets/js/util.js */ str[1366688:1366706],
		/* util.js */ str[1366699:1366706],
		"",
		6520, 1559863020,
		/* application/javascript */ str[1366564:1366586],
		/* S7JxugsJQ2E295zbvpEpk3C59VQ-gz */ str[1365502:1365532],
		true, bytes[61201:63137], str[61201:63137])

	FS.AddFile( /* /assets/sass/base/_page.scss */ str[1366337:1366365],
		/* _page.scss */ str[1366355:1366365],
		"",
		988, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* cmAJ31HROemD7_Vcyvqghs-kb64-gz */ str[1365322:1365352],
		true, bytes[63137:63687], str[63137:63687])

	FS.AddFile( /* /assets/sass/base/_reset.scss */ str[1366222:1366251],
		/* _reset.scss */ str[1366240:1366251],
		"",
		1570, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* _twxRtfqagfQFPRDAmIiD3OhzwQ-gz */ str[1365112:1365142],
		true, bytes[63687:64492], str[63687:64492])

	FS.AddFile( /* /assets/sass/base/_typography.scss */ str[1363537:1363571],
		/* _typography.scss */ str[1363555:1363571],
		"",
		3436, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* d-J3t18Mm3ykHVlPjOdDCAfoGck-gz */ str[1364932:1364962],
		true, bytes[64492:65561], str[64492:65561])

	FS.AddFile( /* /assets/sass/components/_actions.scss */ str[1362837:1362874],
		/* _actions.scss */ str[1362861:1362874],
		"",
		1789, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* vuvB0I6MB6eqgVhIGtpAv2brX7E-gz */ str[1364782:1364812],
		true, bytes[65561:66204], str[65561:66204])

	FS.AddFile( /* /assets/sass/components/_box.scss */ str[1363605:1363638],
		/* _box.scss */ str[1363629:1363638],
		"",
		532, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* hAJb4h0JCvFnE3AMGS-D0IHwbjY-gz */ str[1364572:1364602],
		true, bytes[66204:66499], str[66204:66499])

	FS.AddFile( /* /assets/sass/components/_button.scss */ str[1362911:1362947],
		/* _button.scss */ str[1362935:1362947],
		"",
		2084, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* lxed2eCkzPzpfv1ENtL_l8SHVZ8-gz */ str[1364122:1364152],
		true, bytes[66499:67301], str[66499:67301])

	FS.AddFile( /* /assets/sass/components/_contact.scss */ str[1362763:1362800],
		/* _contact.scss */ str[1362787:1362800],
		"",
		312, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* pr4Ci3KQ3edNmfXeK7HVYSMY5mE-gz */ str[1365922:1365952],
		true, bytes[67301:67538], str[67301:67538])

	FS.AddFile( /* /assets/sass/components/_features.scss */ str[1362688:1362726],
		/* _features.scss */ str[1362712:1362726],
		"",
		2001, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* npXndGZcl0mS_lTskRSXnbthUq4-gz */ str[1364902:1364932],
		true, bytes[67538:68262], str[67538:68262])

	FS.AddFile( /* /assets/sass/components/_form.scss */ str[1363401:1363435],
		/* _form.scss */ str[1363425:1363435],
		"",
		5208, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* vdYhYP5wtdmVflKfBTWMrOykMnQ-gz */ str[1365022:1365052],
		true, bytes[68262:69806], str[68262:69806])

	FS.AddFile( /* /assets/sass/components/_icon.scss */ str[1363333:1363367],
		/* _icon.scss */ str[1363357:1363367],
		"",
		1210, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* 0e-_tzvDdGI4gvY9QDxFmd4km5A-gz */ str[1365712:1365742],
		true, bytes[69806:70307], str[69806:70307])

	FS.AddFile( /* /assets/sass/components/_icons.scss */ str[1362983:1363018],
		/* _icons.scss */ str[1363007:1363018],
		"",
		423, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* 71jDOmCcvPnyuVHrnMT7xjEsi_4-gz */ str[1364542:1364572],
		true, bytes[70307:70602], str[70307:70602])

	FS.AddFile( /* /assets/sass/components/_image.scss */ str[1363298:1363333],
		/* _image.scss */ str[1363322:1363333],
		"",
		886, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* JPidTwOeGVwTei8HzpgFNILKtxo-gz */ str[1365412:1365442],
		true, bytes[70602:70975], str[70602:70975])

	FS.AddFile( /* /assets/sass/components/_list.scss */ str[1363503:1363537],
		/* _list.scss */ str[1363527:1363537],
		"",
		909, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* mtkzYnNTfEhq109WXaXjf5tK7Jc-gz */ str[1364302:1364332],
		true, bytes[70975:71381], str[70975:71381])

	FS.AddFile( /* /assets/sass/components/_menu.scss */ str[1363571:1363605],
		/* _menu.scss */ str[1363595:1363605],
		"",
		609, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* rFvVzEekW7AQjw0WUvPnBNxfjNE-gz */ str[1365952:1365982],
		true, bytes[71381:71733], str[71381:71733])

	FS.AddFile( /* /assets/sass/components/_row.scss */ str[1363638:1363671],
		/* _row.scss */ str[1363662:1363671],
		"",
		631, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* YcdwTsrPZvgElPm1foCNsIN6DjM-gz */ str[1364632:1364662],
		true, bytes[71733:71987], str[71733:71987])

	FS.AddFile( /* /assets/sass/components/_section.scss */ str[1362800:1362837],
		/* _section.scss */ str[1362824:1362837],
		"",
		744, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* R3Tp5DEnmYIGvczunCfXGNzorbc-gz */ str[1364812:1364842],
		true, bytes[71987:72347], str[71987:72347])

	FS.AddFile( /* /assets/sass/components/_split.scss */ str[1363263:1363298],
		/* _split.scss */ str[1363287:1363298],
		"",
		1490, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* BL2vDRwurAq07SsWo9kXZFs76yM-gz */ str[1365202:1365232],
		true, bytes[72347:72847], str[72347:72847])

	FS.AddFile( /* /assets/sass/components/_spotlights.scss */ str[1362648:1362688],
		/* _spotlights.scss */ str[1362672:1362688],
		"",
		2473, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* lCfK-nZwct547sO6M5rHe7LzxrQ-gz */ str[1365382:1365412],
		true, bytes[72847:73651], str[72847:73651])

	FS.AddFile( /* /assets/sass/components/_table.scss */ str[1363018:1363053],
		/* _table.scss */ str[1363042:1363053],
		"",
		1398, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* whjgBpUuzXX9KMOZqPSQIs7uYKg-gz */ str[1366192:1366222],
		true, bytes[73651:74199], str[73651:74199])

	FS.AddFile( /* /assets/sass/components/_wrapper.scss */ str[1362726:1362763],
		/* _wrapper.scss */ str[1362750:1362763],
		"",
		2409, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* 3mtvOygf9cEnasSkd5TNvzVq4Tk-gz */ str[1365772:1365802],
		true, bytes[74199:74849], str[74199:74849])

	FS.AddFile( /* /assets/sass/layout/_footer.scss */ str[1363999:1364031],
		/* _footer.scss */ str[1364019:1364031],
		"",
		618, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* 8SE0UhiPh3ZnW1SIJtf5K5Il6jY-gz */ str[1366072:1366102],
		true, bytes[74849:75203], str[74849:75203])

	FS.AddFile( /* /assets/sass/layout/_header.scss */ str[1363935:1363967],
		/* _header.scss */ str[1363955:1363967],
		"",
		1592, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* YYeyQP2D2f7pwebmKvpJ19pSCQU-gz */ str[1365682:1365712],
		true, bytes[75203:75838], str[75203:75838])

	FS.AddFile( /* /assets/sass/layout/_intro.scss */ str[1364031:1364062],
		/* _intro.scss */ str[1364051:1364062],
		"",
		627, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* uoZlJ9AWiZoOFi57lzFjxavoMqk-gz */ str[1366132:1366162],
		true, bytes[75838:76185], str[75838:76185])

	FS.AddFile( /* /assets/sass/layout/_sidebar.scss */ str[1363803:1363836],
		/* _sidebar.scss */ str[1363823:1363836],
		"",
		3771, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* 90-X22gqfF84OnFqI3Ko531wtwM-gz */ str[1365652:1365682],
		true, bytes[76185:77386], str[76185:77386])

	FS.AddFile( /* /assets/sass/layout/_wrapper.scss */ str[1363836:1363869],
		/* _wrapper.scss */ str[1362750:1362763],
		"",
		522, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* -LlEHMgF4AF29o7taqHiOwuxmwI-gz */ str[1365472:1365502],
		true, bytes[77386:77698], str[77386:77698])

	FS.AddFile( /* /assets/sass/libs/_breakpoints.scss */ str[1363228:1363263],
		/* _breakpoints.scss */ str[1363246:1363263],
		"",
		4577, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* vOtt1jIQU5IfMkimrqfFWBgp4EA-gz */ str[1365352:1365382],
		true, bytes[77698:78652], str[77698:78652])

	FS.AddFile( /* /assets/sass/libs/_functions.scss */ str[1363770:1363803],
		/* _functions.scss */ str[1363788:1363803],
		"",
		1957, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* MQph4TOx2mnmTA1zzENEr_DLHeo-gz */ str[1365052:1365082],
		true, bytes[78652:79259], str[78652:79259])

	FS.AddFile( /* /assets/sass/libs/_html-grid.scss */ str[1363737:1363770],
		/* _html-grid.scss */ str[1363755:1363770],
		"",
		2840, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* maaMKjkDtHlTL3qeQw8cxM91MOw-gz */ str[1366162:1366192],
		true, bytes[79259:80213], str[79259:80213])

	FS.AddFile( /* /assets/sass/libs/_mixins.scss */ str[1364242:1364272],
		/* _mixins.scss */ str[1364260:1364272],
		"",
		2218, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* 8Ycu9Z5TTTV_wy4AlNBEwGbX1XY-gz */ str[1364842:1364872],
		true, bytes[80213:81135], str[80213:81135])

	FS.AddFile( /* /assets/sass/libs/_vars.scss */ str[1366309:1366337],
		/* _vars.scss */ str[1366327:1366337],
		"",
		1040, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* Xb5VWWYSRLobbxFGs2ZUy4LS7fs-gz */ str[1364752:1364782],
		true, bytes[81135:81626], str[81135:81626])

	FS.AddFile( /* /assets/sass/libs/_vendor.scss */ str[1364362:1364392],
		/* _vendor.scss */ str[1364380:1364392],
		"",
		7355, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* CxBy6B46ewbN731c0oarze3Qd84-gz */ str[1364692:1364722],
		true, bytes[81626:83929], str[81626:83929])

	FS.AddFile( /* /assets/sass/main.scss */ str[1366586:1366608],
		/* main.scss */ str[1366599:1366608],
		"",
		1367, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* QvVc5CghjxZSyrnum0ZbstgOHjI-gz */ str[1364602:1364632],
		true, bytes[83929:84431], str[83929:84431])

	FS.AddFile( /* /assets/sass/noscript.scss */ str[1366393:1366419],
		/* noscript.scss */ str[1366406:1366419],
		"",
		916, 1559863020,
		/* text/plain; charset=utf-8 */ str[1366419:1366444],
		/* QawZPTEpGdxALGWu-uXK6IysPj0-gz */ str[1364512:1364542],
		true, bytes[84431:84806], str[84431:84806])

	FS.AddFile( /* /assets/webfonts/fa-brands-400.eot */ str[1363367:1363401],
		/* fa-brands-400.eot */ str[1363384:1363401],
		"",
		129352, 1559863020,
		/* application/vnd.ms-fontobject */ str[1366251:1366280],
		/* ytMxc0Lb454zNWljPxwQPh0ThyI-gz */ str[1364452:1364482],
		true, bytes[84806:172766], str[84806:172766])

	FS.AddFile( /* /assets/webfonts/fa-brands-400.svg */ str[1363435:1363469],
		/* fa-brands-400.svg */ str[1363452:1363469],
		"",
		635379, 1559863020,
		/* image/svg+xml */ str[1366822:1366835],
		/* aZdGWKrzrlpmBXEbVNIQ-euLGno-gz */ str[1364392:1364422],
		true, bytes[172766:392898], str[172766:392898])

	FS.AddFile( /* /assets/webfonts/fa-brands-400.ttf */ str[1363469:1363503],
		/* fa-brands-400.ttf */ str[1363486:1363503],
		"",
		129048, 1559863020,
		/* application/font-sfnt */ str[1366629:1366650],
		/* B9Pxf6ht31jIRC7DusvO6Q_kK9w-gz */ str[1364332:1364362],
		true, bytes[392898:480704], str[392898:480704])

	FS.AddFile( /* /assets/webfonts/fa-brands-400.woff */ str[1362947:1362982],
		/* fa-brands-400.woff */ str[1362964:1362982],
		"",
		87352, 1559863020,
		/* application/font-woff */ str[1366608:1366629],
		/* TLCHA169OvoMWI8SnRVJ5TLAbRg-gz */ str[1364272:1364302],
		false, bytes[480704:568056], str[480704:568056])

	FS.AddFile( /* /assets/webfonts/fa-brands-400.woff2 */ str[1362947:1362983],
		/* fa-brands-400.woff2 */ str[1362964:1362983],
		"",
		74508, 1559863020,
		/* font/woff2 */ str[1366882:1366892],
		/* pAOvMzfmIH0US5mLnDvtQ5r1Yqk-gz */ str[1364212:1364242],
		true, bytes[568056:642557], str[568056:642557])

	FS.AddFile( /* /assets/webfonts/fa-regular-400.eot */ str[1363158:1363193],
		/* fa-regular-400.eot */ str[1363175:1363193],
		"",
		34388, 1559863020,
		/* application/vnd.ms-fontobject */ str[1366251:1366280],
		/* OLBdqwMqFO2QTDaHd5W-l0F8s88-gz */ str[1364092:1364122],
		true, bytes[642557:659591], str[642557:659591])

	FS.AddFile( /* /assets/webfonts/fa-regular-400.svg */ str[1363088:1363123],
		/* fa-regular-400.svg */ str[1363105:1363123],
		"",
		132768, 1559863020,
		/* image/svg+xml */ str[1366822:1366835],
		/* 22nrzGmEnLV3jpIxriG8SKSlNdM-gz */ str[1364662:1364692],
		true, bytes[659591:694752], str[659591:694752])

	FS.AddFile( /* /assets/webfonts/fa-regular-400.ttf */ str[1363123:1363158],
		/* fa-regular-400.ttf */ str[1363140:1363158],
		"",
		34092, 1559863020,
		/* application/font-sfnt */ str[1366629:1366650],
		/* KDP0hsmANCdqjPZj9jm9mF82C9A-gz */ str[1365832:1365862],
		true, bytes[694752:711726], str[694752:711726])

	FS.AddFile( /* /assets/webfonts/fa-regular-400.woff */ str[1362874:1362910],
		/* fa-regular-400.woff */ str[1362891:1362910],
		"",
		16804, 1559863020,
		/* application/font-woff */ str[1366608:1366629],
		/* godbPjH0oqi57U5LOD5hzjc3Jn8-gz */ str[1365982:1366012],
		false, bytes[711726:728530], str[711726:728530])

	FS.AddFile( /* /assets/webfonts/fa-regular-400.woff2 */ str[1362874:1362911],
		/* fa-regular-400.woff2 */ str[1362891:1362911],
		"",
		13580, 1559863020,
		/* font/woff2 */ str[1366882:1366892],
		/* x0QhfKqCsyRc_6JxSq8uyfdJYU0-gz */ str[1364062:1364092],
		false, bytes[728530:742110], str[728530:742110])

	FS.AddFile( /* /assets/webfonts/fa-solid-900.eot */ str[1363704:1363737],
		/* fa-solid-900.eot */ str[1363721:1363737],
		"",
		192116, 1559863020,
		/* application/vnd.ms-fontobject */ str[1366251:1366280],
		/* rO5PHjYcz5lRIfOBuVQktuxrKcA-gz */ str[1365562:1365592],
		true, bytes[742110:841794], str[742110:841794])

	FS.AddFile( /* /assets/webfonts/fa-solid-900.svg */ str[1363902:1363935],
		/* fa-solid-900.svg */ str[1363919:1363935],
		"",
		777689, 1559863020,
		/* image/svg+xml */ str[1366822:1366835],
		/* 2Sswuf1PuuDExya95WE9R1Jixzk-gz */ str[1364482:1364512],
		true, bytes[841794:1062567], str[841794:1062567])

	FS.AddFile( /* /assets/webfonts/fa-solid-900.ttf */ str[1363671:1363704],
		/* fa-solid-900.ttf */ str[1363688:1363704],
		"",
		191832, 1559863020,
		/* application/font-sfnt */ str[1366629:1366650],
		/* t_co7Fke2d2QKO0962pTbQA4yNY-gz */ str[1364182:1364212],
		true, bytes[1062567:1162134], str[1062567:1162134])

	FS.AddFile( /* /assets/webfonts/fa-solid-900.woff */ str[1363053:1363087],
		/* fa-solid-900.woff */ str[1363070:1363087],
		"",
		98020, 1559863020,
		/* application/font-woff */ str[1366608:1366629],
		/* D2_L9H3tGhxqzg9nOD1RK2oFPa0-gz */ str[1364422:1364452],
		false, bytes[1162134:1260154], str[1162134:1260154])

	FS.AddFile( /* /assets/webfonts/fa-solid-900.woff2 */ str[1363053:1363088],
		/* fa-solid-900.woff2 */ str[1363070:1363088],
		"",
		75440, 1559863020,
		/* font/woff2 */ str[1366882:1366892],
		/* B77RU9R_kSmpRO5U3XKVLe7QdMg-gz */ str[1364872:1364902],
		true, bytes[1260154:1335565], str[1260154:1335565])

	FS.AddFile( /* /elements.html */ str[1366808:1366822],
		/* elements.html */ str[1366809:1366822],
		"",
		11933, 1559863020,
		/* text/html; charset=utf-8 */ str[1366469:1366493],
		/* hGCl2dwd9BchHhiPja8yxWxUP-c-gz */ str[1364962:1364992],
		true, bytes[1335565:1337938], str[1335565:1337938])

	FS.AddFile( /* /generic.html */ str[1366835:1366848],
		/* generic.html */ str[1366836:1366848],
		"",
		2159, 1559863020,
		/* text/html; charset=utf-8 */ str[1366469:1366493],
		/* z2TbXtdU9mD79eXCpSkRY7xxUG4-gz */ str[1365082:1365112],
		true, bytes[1337938:1338766], str[1337938:1338766])

	FS.AddFile( /* /images/pic01.jpg */ str[1366740:1366757],
		/* pic01.jpg */ str[1366748:1366757],
		"",
		6953, 1559863020,
		/* image/jpeg */ str[1366892:1366902],
		/* wCYOBX-xnmqOg54pM4OPB23uiLo-gz */ str[1365172:1365202],
		true, bytes[1338766:1343113], str[1338766:1343113])

	FS.AddFile( /* /images/pic02.jpg */ str[1366757:1366774],
		/* pic02.jpg */ str[1366765:1366774],
		"",
		5767, 1559863020,
		/* image/jpeg */ str[1366892:1366902],
		/* nQH7Cczk_fZolUkELQQLh3prJP4-gz */ str[1365232:1365262],
		true, bytes[1343113:1346986], str[1343113:1346986])

	FS.AddFile( /* /images/pic03.jpg */ str[1366723:1366740],
		/* pic03.jpg */ str[1366731:1366740],
		"",
		6828, 1559863020,
		/* image/jpeg */ str[1366892:1366902],
		/* -CeFxiCs0ntnmS8zTsCps3Ku-t8-gz */ str[1365262:1365292],
		true, bytes[1346986:1350768], str[1346986:1350768])

	FS.AddFile( /* /images/pic04.jpg */ str[1366791:1366808],
		/* pic04.jpg */ str[1366799:1366808],
		"",
		12171, 1559863020,
		/* image/jpeg */ str[1366892:1366902],
		/* 02qb4z_4YIIlHjnPZHs011ycohk-gz */ str[1365532:1365562],
		true, bytes[1350768:1357887], str[1350768:1357887])

	FS.AddFile( /* /images/pic05.jpg */ str[1366774:1366791],
		/* pic05.jpg */ str[1366782:1366791],
		"",
		2527, 1559863020,
		/* image/jpeg */ str[1366892:1366902],
		/* -O7aPD8LX1TqYA4koUxZXVrOqUY-gz */ str[1365802:1365832],
		true, bytes[1357887:1359373], str[1357887:1359373])

	FS.AddFile( /* /images/pic06.jpg */ str[1366706:1366723],
		/* pic06.jpg */ str[1366714:1366723],
		"",
		2798, 1559863020,
		/* image/jpeg */ str[1366892:1366902],
		/* NvP1WxPrrMgeikQqL85wTo4iS_0-gz */ str[1366102:1366132],
		true, bytes[1359373:1361165], str[1359373:1361165])

	FS.AddFile( /* /index.html */ str[1366860:1366871],
		/* index.html */ str[1366861:1366871],
		"",
		5954, 1559863020,
		/* text/html; charset=utf-8 */ str[1366469:1366493],
		/* 0JawvkBv9UaGjnaJ5t0elGNSBHs-gz */ str[1364152:1364182],
		true, bytes[1361165:1362648], str[1361165:1362648])

	FS.AddFolder( /* /images */ str[1366376:1366383],
		/* images */ str[1366377:1366383],
		"",
		1559863020,
		/* /images/pic01.jpg */ str[1366740:1366757],
		/* /images/pic02.jpg */ str[1366757:1366774],
		/* /images/pic03.jpg */ str[1366723:1366740],
		/* /images/pic04.jpg */ str[1366791:1366808],
		/* /images/pic05.jpg */ str[1366774:1366791],
		/* /images/pic06.jpg */ str[1366706:1366723],
	)

	FS.AddFolder( /* /assets/webfonts */ str[1362874:1362890],
		/* webfonts */ str[1362882:1362890],
		"",
		1559863020,
		/* /assets/webfonts/fa-brands-400.eot */ str[1363367:1363401],
		/* /assets/webfonts/fa-brands-400.svg */ str[1363435:1363469],
		/* /assets/webfonts/fa-brands-400.ttf */ str[1363469:1363503],
		/* /assets/webfonts/fa-brands-400.woff */ str[1362947:1362982],
		/* /assets/webfonts/fa-brands-400.woff2 */ str[1362947:1362983],
		/* /assets/webfonts/fa-regular-400.eot */ str[1363158:1363193],
		/* /assets/webfonts/fa-regular-400.svg */ str[1363088:1363123],
		/* /assets/webfonts/fa-regular-400.ttf */ str[1363123:1363158],
		/* /assets/webfonts/fa-regular-400.woff */ str[1362874:1362910],
		/* /assets/webfonts/fa-regular-400.woff2 */ str[1362874:1362911],
		/* /assets/webfonts/fa-solid-900.eot */ str[1363704:1363737],
		/* /assets/webfonts/fa-solid-900.svg */ str[1363902:1363935],
		/* /assets/webfonts/fa-solid-900.ttf */ str[1363671:1363704],
		/* /assets/webfonts/fa-solid-900.woff */ str[1363053:1363087],
		/* /assets/webfonts/fa-solid-900.woff2 */ str[1363053:1363088],
	)

	FS.AddFolder( /* /assets/sass/libs */ str[1363228:1363245],
		/* libs */ str[1363241:1363245],
		"",
		1559863020,
		/* /assets/sass/libs/_breakpoints.scss */ str[1363228:1363263],
		/* /assets/sass/libs/_functions.scss */ str[1363770:1363803],
		/* /assets/sass/libs/_html-grid.scss */ str[1363737:1363770],
		/* /assets/sass/libs/_mixins.scss */ str[1364242:1364272],
		/* /assets/sass/libs/_vars.scss */ str[1366309:1366337],
		/* /assets/sass/libs/_vendor.scss */ str[1364362:1364392],
	)

	FS.AddFolder( /* /assets/sass/layout */ str[1363803:1363822],
		/* layout */ str[1363816:1363822],
		"",
		1559863020,
		/* /assets/sass/layout/_footer.scss */ str[1363999:1364031],
		/* /assets/sass/layout/_header.scss */ str[1363935:1363967],
		/* /assets/sass/layout/_intro.scss */ str[1364031:1364062],
		/* /assets/sass/layout/_sidebar.scss */ str[1363803:1363836],
		/* /assets/sass/layout/_wrapper.scss */ str[1363836:1363869],
	)

	FS.AddFolder( /* /assets/sass/components */ str[1362648:1362671],
		/* components */ str[1362661:1362671],
		"",
		1559863020,
		/* /assets/sass/components/_actions.scss */ str[1362837:1362874],
		/* /assets/sass/components/_box.scss */ str[1363605:1363638],
		/* /assets/sass/components/_button.scss */ str[1362911:1362947],
		/* /assets/sass/components/_contact.scss */ str[1362763:1362800],
		/* /assets/sass/components/_features.scss */ str[1362688:1362726],
		/* /assets/sass/components/_form.scss */ str[1363401:1363435],
		/* /assets/sass/components/_icon.scss */ str[1363333:1363367],
		/* /assets/sass/components/_icons.scss */ str[1362983:1363018],
		/* /assets/sass/components/_image.scss */ str[1363298:1363333],
		/* /assets/sass/components/_list.scss */ str[1363503:1363537],
		/* /assets/sass/components/_menu.scss */ str[1363571:1363605],
		/* /assets/sass/components/_row.scss */ str[1363638:1363671],
		/* /assets/sass/components/_section.scss */ str[1362800:1362837],
		/* /assets/sass/components/_split.scss */ str[1363263:1363298],
		/* /assets/sass/components/_spotlights.scss */ str[1362648:1362688],
		/* /assets/sass/components/_table.scss */ str[1363018:1363053],
		/* /assets/sass/components/_wrapper.scss */ str[1362726:1362763],
	)

	FS.AddFolder( /* /assets/sass/base */ str[1363537:1363554],
		/* base */ str[1363550:1363554],
		"",
		1559863020,
		/* /assets/sass/base/_page.scss */ str[1366337:1366365],
		/* /assets/sass/base/_reset.scss */ str[1366222:1366251],
		/* /assets/sass/base/_typography.scss */ str[1363537:1363571],
	)

	FS.AddFolder( /* /assets/sass */ str[1362648:1362660],
		/* sass */ str[1362656:1362660],
		"",
		1559863020,
		/* /assets/sass/base */ str[1363537:1363554],
		/* /assets/sass/components */ str[1362648:1362671],
		/* /assets/sass/layout */ str[1363803:1363822],
		/* /assets/sass/libs */ str[1363228:1363245],
		/* /assets/sass/main.scss */ str[1366586:1366608],
		/* /assets/sass/noscript.scss */ str[1366393:1366419],
	)

	FS.AddFolder( /* /assets/js */ str[1363869:1363879],
		/* js */ str[1363877:1363879],
		"",
		1559863020,
		/* /assets/js/breakpoints.min.js */ str[1366280:1366309],
		/* /assets/js/browser.min.js */ str[1366444:1366469],
		/* /assets/js/jquery.min.js */ str[1366493:1366517],
		/* /assets/js/jquery.scrollex.min.js */ str[1363869:1363902],
		/* /assets/js/jquery.scrolly.min.js */ str[1363967:1363999],
		/* /assets/js/main.js */ str[1366670:1366688],
		/* /assets/js/util.js */ str[1366688:1366706],
	)

	FS.AddFolder( /* /assets/css/images */ str[1366365:1366383],
		/* images */ str[1366377:1366383],
		"",
		1559863020,
		/* /assets/css/images/intro.svg */ str[1366365:1366393],
	)

	FS.AddFolder( /* /assets/css */ str[1363193:1363204],
		/* css */ str[1362685:1362688],
		"",
		1559863020,
		/* /assets/css/fontawesome-all.min.css */ str[1363193:1363228],
		/* /assets/css/images */ str[1366365:1366383],
		/* /assets/css/main.css */ str[1366650:1366670],
		/* /assets/css/noscript.css */ str[1366517:1366541],
	)

	FS.AddFolder( /* /assets */ str[1362648:1362655],
		/* assets */ str[1362649:1362655],
		"",
		1559863020,
		/* /assets/css */ str[1363193:1363204],
		/* /assets/js */ str[1363869:1363879],
		/* /assets/sass */ str[1362648:1362660],
		/* /assets/webfonts */ str[1362874:1362890],
	)

	FS.AddFolder( /* / */ str[1362648:1362649],
		/* / */ str[1362648:1362649],
		"",
		1581646592,
		/* /LICENSE.txt */ str[1366848:1366860],
		/* /README.txt */ str[1366871:1366882],
		/* /assets */ str[1362648:1362655],
		/* /elements.html */ str[1366808:1366822],
		/* /generic.html */ str[1366835:1366848],
		/* /images */ str[1366376:1366383],
		/* /index.html */ str[1366860:1366871],
	)
}
