package days

import (
	"fmt"
	"sync"

	u "parnic.com/aoc2019/utilities"
)

type day23Computer struct {
	program            u.IntcodeProgram
	id                 int
	packetQueue        chan u.Vec2[int64]
	outputStep         int
	nextPacketDest     int
	sendingPacket      u.Vec2[int64]
	hasQueuedPacket    bool
	lastReceivedPacket u.Vec2[int64]
	idle               bool
}

func (d Day23) makeComputer(id int) *day23Computer {
	c := &day23Computer{
		program:     d.program.Copy(),
		id:          id,
		packetQueue: make(chan u.Vec2[int64]),
		idle:        true,
	}

	return c
}

type Day23 struct {
	program u.IntcodeProgram
}

func (d *Day23) Parse() {
	d.program = u.LoadIntcodeProgram("23p")
}

func (d Day23) Num() int {
	return 23
}

func (d Day23) initComputers() []*day23Computer {
	computers := make([]*day23Computer, 50)
	for i := range computers {
		computers[i] = d.makeComputer(i)
	}
	return computers
}

func (d Day23) execComputers(computers []*day23Computer, nat chan u.Vec2[int64]) *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(len(computers))
	for _, c := range computers {
		go func(c *day23Computer) {
			bootedUp := false
			c.program.RunIn(func(inputStep int) int64 {
				if !bootedUp {
					bootedUp = true
					return int64(c.id)
				}

				if c.hasQueuedPacket {
					// fmt.Printf("  %d finished processing packet %v\n", c.id, c.lastReceivedPacket)
					c.hasQueuedPacket = false
					return c.lastReceivedPacket.Y
				}

				select {
				case c.lastReceivedPacket = <-c.packetQueue:
					// fmt.Printf("computer %d received packet %v\n", c.id, packet)
					c.hasQueuedPacket = true
					return c.lastReceivedPacket.X
				default:
					c.idle = true
					return -1
				}
			}, func(val int64, state u.IntcodeProgramState) {
				c.idle = false
				switch c.outputStep {
				case 0:
					c.nextPacketDest = int(val)
				case 1:
					c.sendingPacket.X = val
				case 2:
					c.sendingPacket.Y = val

					if c.nextPacketDest == 255 {
						// fmt.Printf("computer %d sending %v to 255\n", c.id, c.sendingPacket)
						nat <- c.sendingPacket
					} else {
						// fmt.Printf("computer %d sending %v to computer %d\n", c.id, c.sendingPacket, c.nextPacketDest)
						computers[c.nextPacketDest].packetQueue <- c.sendingPacket
					}
				}

				c.outputStep = (c.outputStep + 1) % 3
			})

			wg.Done()
		}(c)
	}

	return wg
}

func (d *Day23) Part1() string {
	computers := d.initComputers()
	natChan := make(chan u.Vec2[int64])
	wg := d.execComputers(computers, natChan)

	answer := <-natChan
	for _, c := range computers {
		c.program.Stop()
	}
	// not really necessary, but let's make sure they all shut down in case
	// we're running all days at once
	wg.Wait()

	return fmt.Sprintf("First packet sent to 255 Y value: %s%d%s", u.TextBold, answer.Y, u.TextReset)
}

func (d *Day23) Part2() string {
	computers := d.initComputers()
	natChan := make(chan u.Vec2[int64])
	wg := d.execComputers(computers, natChan)

	answerChan := make(chan int64)

	go func() {
		var currVal u.Vec2[int64]
		var lastVal u.Vec2[int64]
		hasReceived := false
		for {
			select {
			case currVal = <-natChan:
				hasReceived = true
			default:
			}

			allIdle := true
			for _, c := range computers {
				if !c.idle {
					allIdle = false
					break
				}
			}

			if allIdle && hasReceived {
				// fmt.Printf("all idle, sending %v to computer 0\n", currVal)
				if lastVal.Y == currVal.Y {
					// fmt.Printf("found answer? %d\n", currVal.Y)
					answerChan <- currVal.Y
				}
				computers[0].packetQueue <- currVal
				lastVal = currVal
			}
		}
	}()

	answer := <-answerChan
	for _, c := range computers {
		c.program.Stop()
	}
	wg.Wait()

	return fmt.Sprintf("First Y value sent to the NAT twice in a row: %s%d%s", u.TextBold, answer, u.TextReset)
}
