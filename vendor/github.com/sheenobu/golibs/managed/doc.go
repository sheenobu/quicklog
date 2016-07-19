package managed

/*
	system := managed.System("name")

	system.Add(managed.Simple("name", func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.After(1 * time.Second):

			}
		}
	}))

	system.Add(managed.Timer("name", 4 * time.Second, func(ctx context.Context, t time.Time) {
		// do X
	}))

	system.Add(managed.Nats("name", "manager", func(ctx context.Context, msg []byte) {

	}))

	system.Require("name1", "name2")

	system.Get("name").Stop()


	system.Run()
*/
