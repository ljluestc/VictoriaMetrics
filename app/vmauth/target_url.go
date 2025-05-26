func (ui *UserInfo) getURLPrefixAndHeaders(u *url.URL, host string, h http.Header) (string, map[string]string) {
	for _, e := range ui.URLMaps {
		if !matchAnyRegex(e.SrcHosts, host) {
			continue
		}
		if !matchAnyRegex(e.SrcPaths, u.Path) {
			continue
		}
		if !matchAnyQueryArg(e.SrcQueryArgs, u.Query()) {
			continue
		}
		if !matchAnyHeader(e.SrcHeaders, h) {
			continue
		}
		return e.URLPrefix.URL, convertHeadersConfToMap(e.HeadersConf)
	}
	if ui.URLPrefix != nil {
		return ui.URLPrefix.URL, nil
	}
	return "", nil
}
