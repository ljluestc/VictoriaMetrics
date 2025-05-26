// ...existing code...

func (c *Conn) readInitialHandshake() error {
	// ...existing code...

	if len(data) > pos {
		// ...existing code...

		if c.capability&mysql.CLIENT_SECURE_CONNECTION != 0 {
			// ...existing code...

			if pos+rest-1 > len(data) {
				return errors.New("malformed packet: insufficient data for auth plugin scramble")
			}
			authPluginDataPart2 := data[pos : pos+rest-1]
			pos += rest

			c.salt = append(c.salt, authPluginDataPart2...)
		}

		if c.capability&mysql.CLIENT_PLUGIN_AUTH != 0 {
			if pos >= len(data) {
				return errors.New("malformed packet: missing auth plugin name")
			}
			authPluginNameEnd := bytes.IndexByte(data[pos:], 0x00)
			if authPluginNameEnd == -1 {
				return errors.New("malformed packet: auth plugin name not null-terminated")
			}
			c.authPluginName = string(data[pos : pos+authPluginNameEnd])
			pos += authPluginNameEnd + 1

			if pos < len(data) && data[pos] != 0 {
				return errors.Errorf("expect 0x00 after authPluginName, got %q", rune(data[pos]))
			}
		}
	}

	// ...existing code...
	return nil
}

func (c *Conn) writeAuthHandshake() error {
	// ...existing code...

	// Validate collation before proceeding
	collationName := c.collation
	if len(collationName) == 0 {
		collationName = mysql.DEFAULT_COLLATION_NAME
	}
	collation, err := charset.GetCollationByName(collationName)
	if err != nil {
		return fmt.Errorf("invalid collation name %s", collationName)
	}

	// ...existing code...

	// Ensure TLS handshake is successful
	if c.tlsConfig != nil {
		// ...existing code...

		if err := tlsConn.Handshake(); err != nil {
			return fmt.Errorf("TLS handshake failed: %w", err)
		}

		// ...existing code...
	}

	// ...existing code...

	// Validate final packet length
	if pos > len(data) {
		return errors.New("malformed packet: data length mismatch")
	}

	return c.WritePacket(data)
}
