package device

import (
	"fmt"
	"log"
	"sync"

	"github.com/jochenvg/go-udev"
	"github.com/tarm/serial"
)

// Manager manages devices that are plugged into the system. It supports auto
// detection of devices.
//
// Serial ports are opened each for a device, and a clean API for communicating
// is provided via Read, Write and Flush methods.
//
// The devices are monitored via udev, and any changes that requires reloading
// of the  ports are handled by reloading the ports to the devices.
//
// This is safe to use concurrently in multiple goroutines
type Manager struct {
	devices map[string]serial.Config
	conn    []*Conn
	mu      sync.RWMutex
	monitor *udev.Monitor
	done    chan struct{}
	stop    chan struct{}
}

// New returns a new Manager instance
func New() *Manager {
	return &Manager{
		devices: make(map[string]serial.Config),
		done:    make(chan struct{}),
		stop:    make(chan struct{}),
	}
}

// Init initializes the manager. This involves creating a new goroutine to watch
// over the changes detected by udev for any device interaction with the system.
//
// The only interesting device actions are add and reomove for adding and
// removing devices respctively.
func (m *Manager) Init() {
	u := udev.Udev{}
	monitor := u.NewMonitorFromNetlink("udev")
	devCh, err := monitor.DeviceChan(m.done)
	if err != nil {
		panic(err)
	}
	m.monitor = monitor
	go func() {
	stop:
		for {
			select {
			case d := <-devCh:
				switch d.Action() {
				case "add":
					fmt.Printf(" new device added ad %s\n", d.Devpath())
					m.AddDevice(d.Devpath())
				case "remove":
					fmt.Printf(" %s was removed\n", d.Devpath())
					m.RemoveDevice(d.Devpath())
				default:
					fmt.Println(d.Action())
				}
			case quit := <-m.stop:
				m.done <- quit
				break stop
			}
		}
	}()
}

// AddDevice adds device name to the manager
func (m *Manager) AddDevice(name string) error {
	cfg := serial.Config{Name: name}
	m.mu.Lock()
	m.devices[name] = cfg
	m.mu.Unlock()
	return nil
}

// RemoveDevice removes device name from the manager
func (m *Manager) RemoveDevice(name string) error {
	m.mu.RLock()
	delete(m.devices, name)
	m.mu.RUnlock()
	return nil
}

// close all ports that are open for the devices
func (m *Manager) releaseAllPorts() {
	for _, c := range m.conn {
		err := c.Close()
		if err != nil {
			log.Printf("[ERR] closing port %s %v\n", c.device.Name, err)
		}
	}
}

func (m *Manager) reload() {
	m.releaseAllPorts()
	var conns []*Conn
	for _, v := range m.devises {
		conn := &Conn{device: v}
		err := conn.Open()
		if err != nil {
			log.Printf("[ERR] closing port %s %v\n", c.device.Name, err)
		}
	}
	m.conn = conns
}

//Close shuts down the device manager. This makes sure the udev monitor is
//closed and all goroutines are properly exited.
func (m *Manager) Close() {
	m.stop <- struct{}{}
}

// Conn is a device serial connection
type Conn struct {
	device serial.Config
	port   *serial.Port
	isOpen bool
}

// Open opens a serial port to the undelying device
func (c *Conn) Open() error {
	p, err := serial.OpenPort(&c.device)
	if err != nil {
		return nil
	}
	c.port = p
	c.isOpen = true
	return nil
}

// Close closes the port helt by *Conn.
func (c *Conn) Close() error {
	if c.isOpen {
		return c.port.Close()
	}
	return nil
}