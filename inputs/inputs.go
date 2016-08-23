package inputs

import (
	_ "github.com/sheenobu/quicklog/inputs/nats"   // auto-import
	_ "github.com/sheenobu/quicklog/inputs/stdin"  // auto-import
	_ "github.com/sheenobu/quicklog/inputs/syslog" // auto-import
	_ "github.com/sheenobu/quicklog/inputs/tcp"    // auto-import
	_ "github.com/sheenobu/quicklog/inputs/udp"    // auto-import
)
