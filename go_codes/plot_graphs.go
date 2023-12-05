package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Assuming your model is something like this
type RegistrationAccept struct {
	Timestamp           time.Time `json:"timestamp"`
	Delay               int       `json:"delay"`
	Jitter              float64   `json:"jitter"`
	PacketDeliveryRatio float64   `json:"packetDeliveryRatio"`
}

const DataCollectionAmfRegistration = "datacollection.amf.Registration"

func main() {
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("your_mongo_uri"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// Access collection
	collection := client.Database("nwdaf").Collection(DataCollectionAmfRegistration)

	// Retrieve data from MongoDB
	cursor, err := collection.Find(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	var data []RegistrationAccept
	if err := cursor.All(context.TODO(), &data); err != nil {
		log.Fatal(err)
	}

	// Separate slices for each metric
	var timestamps []time.Time
	var delays []int
	var jitters []float64
	var packetDeliveryRatios []float64

	// Fill the slices
	for _, entry := range data {
		timestamps = append(timestamps, entry.Timestamp)
		delays = append(delays, entry.Delay)
		jitters = append(jitters, entry.Jitter)
		packetDeliveryRatios = append(packetDeliveryRatios, entry.PacketDeliveryRatio)
	}

	// Print or use the obtained slices as needed
	fmt.Println("Timestamps:", timestamps)
	fmt.Println("Delays:", delays)
	fmt.Println("Jitters:", jitters)
	fmt.Println("Packet Delivery Ratios:", packetDeliveryRatios)

	// Plot the delay over time
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}

	delayPoints := make(plotter.XYs, len(data))
	for i, entry := range data {
		delayPoints[i].X = float64(entry.Timestamp.Unix())
		delayPoints[i].Y = float64(entry.Delay)
	}

	delayLine, delayPoints, err := plotter.NewLinePoints(delayPoints)
	if err != nil {
		log.Fatal(err)
	}

	p.Add(delayLine)

	// Save the delay plot to a file
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "delay_plot.png"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Delay plot saved to delay_plot.png")
}
