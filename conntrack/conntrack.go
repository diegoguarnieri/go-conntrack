package conntrack

import (
   "os/exec"
   "bytes"
   "io"
   "strings"
)

type Conntrack struct {
   path     string
   version  string
}

func New() (*Conntrack, error) {
   path, err := exec.LookPath("conntrack")
   if err != nil {
      return nil, err
   }
   
   version, err := getConntrackVersion(path)
   if err != nil {
      return nil, err
   }

   conntrack := Conntrack {
      path: path,
      version: version,
   }

   return &conntrack, nil
}

//conntrack -D -s 192.168.0.1
func (conntrack *Conntrack) DeleteConnectionBySrcIp(ip string) error {
   args := []string{"-D", "-s", ip}

   result, err := conntrack.execute(args)
   if err != nil {
      if strings.Contains(result[0], "0 flow") {
         return nil
      } else {
         return err
      }
   }

   return nil
}

func (conntrack *Conntrack) ListConnections() ([]string, error) {
   args := []string{"-L"}

   result, err := conntrack.execute(args)
   if err != nil {
      return nil, err
   }
   
   var connections []string
   for _, val := range result {
      //fmt.Println("0:", strings.Fields(val)[0])
      //fmt.Println("1:", strings.Fields(val)[1])
      connections = append(connections, val)
   }
   return connections, nil
}

func (conntrack *Conntrack) execute(args []string) ([]string, error) {
   var stdout bytes.Buffer
   var stderr bytes.Buffer

   err := conntrack.runWithOutput(args, &stdout, &stderr)
   if err != nil {
      return []string{stderr.String()}, err
   }
   
   result := strings.Split(stdout.String(), "\n")
   if len(result) > 0 && result[len(result)-1] == "" {
      result = result[:len(result)-1]
   }

   return result, nil
}

func (conntrack *Conntrack) runWithOutput(args []string, stdout io.Writer, stderr io.Writer) error {
   args = append([]string{conntrack.path}, args...)

   //var stderr bytes.Buffer
   cmd := exec.Cmd {
      Path:   conntrack.path,
      Args:   args,
      Stdout: stdout,
      Stderr: stderr,
   }

   err := cmd.Run()
   if err != nil {
      return err
   }

   return nil
}

func (conntrack *Conntrack) Version() string {
   return conntrack.version
}

func getConntrackVersion(path string) (string, error) {
   cmd := exec.Command(path, "--version")
   var out bytes.Buffer
   cmd.Stdout = &out
   
   err := cmd.Run()
   if err != nil {
      return "", err
   }

   return out.String(), nil
}