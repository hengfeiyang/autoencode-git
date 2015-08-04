package util

// CSS压缩处理
func CssCompress(in []byte) ([]byte, error) {
	var (
		i   = 0
		j   = 0
		n   = len(in)
		c   byte
		out []byte
	)
	out = make([]byte, n)
	// 在结尾插入一个空格，防止后面+1判断时溢出
	in = append(in, ' ')
	for i = 0; i < n; i++ {
		c = in[i]
		// 注释处理
		if c == '/' {
			if in[i+1] == '*' {
				// 这里是注释
				// 开始寻找注释的结尾
				i++
				for {
					i++
					if i >= n || (in[i] == '*' && in[i+1] == '/') {
						i++
						break
					}
				}
				continue
			}
		}
		// 换行处理
		if c == '\n' {
			continue
		}
		// 处理tab
		if c == '\t' {
			c = ' '
		}
		// 干掉第一个空格
		if j == 0 && c == ' ' {
			continue
		}
		// 处理,:;后面的空格，同时处理连续空格问题
		if c == ' ' || c == ',' || c == ':' || c == ';' || c == '{' || c == '}' {
			// 后面所有的空格都不要啦
			for {
				i++
				if i >= n || (in[i] != ' ' && in[i] != '\n' && in[i] != '\t') {
					i--
					break
				}
			}
		}
		// 处理{前面的空格
		if c == '{' && out[j-1] == ' ' {
			j--
		}
		// 处理}前面的空格和分号
		if c == '}' && (out[j-1] == ' ' || out[j-1] == ';') {
			j--
		}
		out[j] = c
		j++
	}
	return out[0:j], nil
}
