package conntrack

import (
   "testing"
)

func TestConntrack(t *testing.T) {
   conntrack, err := New()
   if err != nil {
      t.Fatalf("New failed: %v", err)
   }

   //show conntrack version
   t.Logf("Version: %v", conntrack.Version())

   //show actual connections
   connections, err := conntrack.ListConnections()
   if err != nil {
      t.Fatal("ListConnections failed: %v\n", err)
   }
   
   for k, v := range connections {
      t.Logf("[%v] %v\n", k, v)
   }

   //delete connections by source ip address
   err = conntrack.DeleteConnectionBySrcIp("192.168.0.1")
   if err != nil {
      t.Logf("DeleteConnectionBySrcIp failed: %v\n", err)
   }
}
