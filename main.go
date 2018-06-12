package main

func main() {
	account := "GCKX3XVTPVNFXQWLQCIBZX6OOPOIUT7FOAZVNOFCNEIXEZFRFSPNZKZT"
	// Check that the key is valid before doing anything else
	payments := paymentsForAccount(account)
	data := fillInVolumePerPayment(payments)
	orderedData := orderData(data)

	createCSV(orderedData, "result.csv")
}
