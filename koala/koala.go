package koala                                                                                                      
                                                                                                                   
import (                                                                                                           
    "log"
    "time"
    "net/url"                                                                                                      
    "strconv"
    "github.com/gliderlabs/registrator/bridge"                                                                     
    koala "github.com/gtato/Koala/koala-go"                                                                        
)                                                                                                                  

func init() {
    f := new(Factory)
    bridge.Register(f, "koala")
}

type Factory struct{}

func (f *Factory) New(uri *url.URL) bridge.RegistryAdapter {
    koaClient := &koala.Client{Url:"http://"+uri.Host}
    return &KoalaAdapter{client: koaClient}
}

type KoalaAdapter struct {
    client *koala.Client
}


func (r *KoalaAdapter) Ping() error {
    log.Printf("Pinging Koala")
    trials := 5
    var err error
    for trials > 0 {
        err = r.client.Version()
        if err != nil {
            log.Println("%s", err)
        }else{
            return nil
        }
        trials -= 1
        time.Sleep(time.Duration(500)*time.Millisecond)
    } 

    return err
}
 
func (r *KoalaAdapter) Register(service *bridge.Service) error {
    log.Println("Register")
    log.Println("service ip: %s", service.IP)
    port := strconv.Itoa(service.Port)
    params := `{"name":"`+service.Name+`", "host":"` + service.IP + `","port":` + port + `}`
    err := r.client.ApiPost("register",params)
    if err != nil {
        log.Println("koala: failed to register service:", err)
    }
    return err
}

func (r *KoalaAdapter) Deregister(service *bridge.Service) error {
    log.Println("Deregister")
    port := strconv.Itoa(service.Port)
    params := `{"name":"`+service.Name+`", "host":"` + service.IP + `","port":` + port + `}`
    err := r.client.ApiPost("deregister",params)
    if err != nil {
        log.Println("koala: failed to deregister service:", err)
    }
    return err
}                                                                                           
                                                      
func (r *KoalaAdapter) Refresh(service *bridge.Service) error {                    
    return nil                                                  
}
                                                
func (r *KoalaAdapter) Services() ([]*bridge.Service, error) {
    return []*bridge.Service{}, nil  
}         