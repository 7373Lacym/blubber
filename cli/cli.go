package cli
import ("github.com/docker/docker/client")
var Cli, err = client.NewClientWithOpts(client.WithVersion("1.37"))