func getTransaction(c *fiber.Ctx) error {
	var transactions []Transaction
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	for cursor.Next(context.Background()) {
		var transaction Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return err
		}
		transactions = append(transactions, transaction)
	}
	return c.JSON(transactions)
}

func createTransaction(c *fiber.Ctx) error {
	transaction := new(Transaction)

	err := c.BodyParser(transaction)

	if err != nil {
		return err
	}
	insertResult, err := collection.InsertOne(context.Background(), transaction)
	if err != nil {
		return err
	}
	transaction.ID = insertResult.InsertedID.(primitive.ObjectID)
	return c.Status(201).JSON(transaction)

}
