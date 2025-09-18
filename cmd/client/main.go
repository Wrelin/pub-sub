package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}

	gameState := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("%s.%s", routing.PauseKey, username),
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gameState),
	)

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username),
		fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, "*"),
		pubsub.SimpleQueueTransient,
		handlerMove(gameState),
	)

	if err != nil {
		log.Fatalf("could not subscribe json: %v", err)
	}

	for {
		input := gamelogic.GetInput()
		if len(input) <= 0 {
			continue
		}

		switch input[0] {
		case "spawn":
			fmt.Println("Spawn unit")

			err = gameState.CommandSpawn(input)
			if err != nil {
				log.Printf("could not spawn unit: %v", err)
			}
		case "move":
			fmt.Println("Moving army")

			move, err := gameState.CommandMove(input)
			if err != nil {
				log.Printf("could not move army: %v", err)
			}

			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilTopic,
				fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username),
				move,
			)
			if err != nil {
				log.Printf("could not publish game state: %v", err)
			}

			fmt.Println("Successfully moved")
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "status":
			fmt.Println("Game status:")
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("unknown command")
		}
	}
}
