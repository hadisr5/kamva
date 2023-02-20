package main

import (
"context"
"errors"
"log"
"os"
"os/signal"
"sync"
"syscall"
"time"
)

var totalInventory int
var mu sync.Mutex

func AddToInventory(ctx context.Context, inventory int, wg *sync.WaitGroup) {
defer wg.Done()

if err := process(inventory); err != nil {
log.Print("Inventory is not valid")
return
}

select {
case <-ctx.Done():
log.Print("Context is canceled")
return
default:
}

mu.Lock()
totalInventory += inventory
mu.Unlock()
}

func process(inventory int) error {
if inventory < 0 {
return errors.New("inventory could not be negative")
}

// It takes two seconds to process the request.
time.Sleep(time.Second * 2)

log.Print("inventory proceed successfully")
return nil
}

func main() {
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Handle interrupt signal to cancel all processes
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)
go func() {
<-c
log.Printf("Received %v, cancelling processes", c)
cancel()
}()

var wg sync.WaitGroup
for i := 0; i < 10; i++ {
wg.Add(1)
go AddToInventory(ctx, i, &wg)
}
wg.Wait()

log.Printf("Total inventory: %d", totalInventory)
}