package localcommand

import (
	"syscall"
	"time"
        "strconv"
	"os"
	"fmt"

	"github.com/yudai/gotty/server"
//        "github.com/davecgh/go-spew/spew"
)

type Options struct {
	CloseSignal  int `hcl:"close_signal" flagName:"close-signal" flagSName:"" flagDescribe:"Signal sent to the command process when gotty close it (default: SIGHUP)" default:"1"`
	CloseTimeout int `hcl:"close_timeout" flagName:"close-timeout" flagSName:"" flagDescribe:"Time in seconds to force kill process after client is disconnected (default: -1)" default:"-1"`
}

type Factory struct {
	command string
	argv    []string
	options *Options
	opts    []Option
        counter int
}

func NewFactory(command string, argv []string, options *Options) (*Factory, error) {
	opts := []Option{WithCloseSignal(syscall.Signal(options.CloseSignal))}
	if options.CloseTimeout >= 0 {
		opts = append(opts, WithCloseTimeout(time.Duration(options.CloseTimeout)*time.Second))
	}

	return &Factory{
		command: command,
		argv:    argv,
		options: options,
		opts:    opts,
                counter: 0,
	}, nil
}

func (factory *Factory) Name() string {
	return "local command"
}

func (factory *Factory) New(params map[string][]string) (server.Slave, error) {
	argv := make([]string, len(factory.argv))
	copy(argv, factory.argv)
	if params["arg"] != nil && len(params["arg"]) > 0 {
		argv = append(argv, params["arg"]...)
	}


	z := 0
        for _, narg := range argv {
            //NS DEBUG log.Printf( "NS-factory> %s", narg)
            if narg == "%count%" {

	      // NS 16/3 find the next free counter by check if a file by the name /tmp/%d.b64 exists?
	        icnt := 0 
	        for icnt = 0; icnt < 4000; icnt++ {
	            fname := fmt.Sprintf( "/tmp/%d.b64", icnt)
	            _, err := os.Stat(fname)
	            if os.IsNotExist(err) {
	               break
	            }
	        }
                argv[z] = strconv.Itoa( /*factory.counter*/icnt)
            }
            z = z + 1
        }

        // NS DEBUG log.Printf("NS-factory> Running cmd: %s, options: %v", factory.command, factory.opts)
        // NS DEBUG spew.Dump( argv)

        factory.counter = /*factory.counter*/ 1

	return New(factory.command, argv, factory.counter, factory.opts...)
}


