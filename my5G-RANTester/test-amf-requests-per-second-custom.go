package templates

import (
        "math/rand"
        "strconv"
        "sync"
        "time"
)

import (
        log "github.com/sirupsen/logrus"
        "my5G-RANTester/config"
        "my5G-RANTester/internal/control_test_engine/gnb"
        "my5G-RANTester/internal/monitoring"
)

func TestRqsLoopCustom(numRqs int, interval int) int64 {
        wg := sync.WaitGroup{}
        monitor := monitoring.Monitor{
                RqsL: 0,
                RqsG: 0,
        }

        cfg, err := config.GetConfig()
        if err != nil {
                log.Fatal("Error in get configuration")
        }

        for y := 1; y <= interval; y++ {
                monitor.InitRqsLocal()

                // Seed the random number generator with a unique value for each iteration
                rand.Seed(time.Now().UnixNano() + int64(y))

                for i := 1; i <= numRqs; i++ {
                        cfg.GNodeB.PlmnList.GnbId = gnbIdGenerator(i)
                        cfg.GNodeB.ControlIF.Port = getRandomPort()

                        go gnb.InitGnbForLoadSeconds(cfg, &wg, &monitor)
                        wg.Add(1)
                }

                wg.Wait()

                log.Warn("[TESTER][GNB] AMF Responses per Second:", monitor.GetRqsLocal())
                monitor.SetRqsGlobal(monitor.GetRqsLocal())
        }

        return monitor.GetRqsGlobal()
}

func gnbIdGeneratorCustom(i int) string {
        var base string
        switch true {
        case i < 10:
                base = "00000"
        case i < 100:
                base = "0000"
        case i >= 100:
                base = "000"
        }

        gnbId := base + strconv.Itoa(i)
        return gnbId
}

func getRandomPort() int {
        return rand.Intn(10) + 1000
}
