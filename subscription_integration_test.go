package braintree

import (
	"fmt"
	"testing"
	"time"

	"github.com/lionelbarrow/braintree-go/nullable"
)

// This test will fail unless you set up your Braintree sandbox account correctly. See TESTING.md for details.
func TestSubscriptionSimple(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(customer)

	g := testGateway.Subscription()

	// Create
	sub, err := g.Create(&SubscriptionRequest{
		PaymentMethodToken: paymentMethod.GetToken(),
		PlanId:             "test_plan",
	})

	t.Log("sub1", sub)

	if err != nil {
		t.Fatal(err)
	}
	if sub.Id == "" {
		t.Fatal("invalid subscription id")
	}

	// Update
	sub2, err := g.Update(&SubscriptionRequest{
		Id:     sub.Id,
		PlanId: "test_plan_2",
		Options: &SubscriptionOptions{
			ProrateCharges:                       true,
			RevertSubscriptionOnProrationFailure: true,
			StartImmediately:                     true,
		},
	})

	t.Log("sub2", sub2)

	if err != nil {
		t.Fatal(err)
	}
	if sub2.Id != sub.Id {
		t.Fatal(sub2.Id)
	}
	if x := sub2.PlanId; x != "test_plan_2" {
		t.Fatal(x)
	}

	// Find
	sub3, err := g.Find(sub.Id)
	if err != nil {
		t.Fatal(err)
	}
	if sub3.Id != sub2.Id {
		t.Fatal(sub3.Id)
	}

	// Cancel
	_, err = g.Cancel(sub2.Id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSubscriptionAllFieldsWithBillingDayOfMonth(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(customer)

	g := testGateway.Subscription()

	// Create
	billingDayOfMonth := nullable.NewNullInt64(15, true)
	numberOfBillingCycles := nullable.NewNullInt64(2, true)
	sub1, err := g.Create(&SubscriptionRequest{
		PaymentMethodToken:    paymentMethod.GetToken(),
		PlanId:                "test_plan",
		MerchantAccountId:     testMerchantAccountId,
		BillingDayOfMonth:     &billingDayOfMonth,
		NumberOfBillingCycles: &numberOfBillingCycles,
		Price: NewDecimal(100, 2),
		Descriptor: &Descriptor{
			Name:  "Company Name*Product 1",
			Phone: "0000000000",
			URL:   "example.com",
		},
	})

	t.Log("sub1", sub1)

	if err != nil {
		t.Fatal(err)
	}
	if sub1.Id == "" {
		t.Fatal("invalid subscription id")
	}
	if sub1.BillingDayOfMonth != "15" {
		t.Fatalf("got billing day of month %#v, want %#v", sub1.BillingDayOfMonth, "15")
	}
	if x := sub1.NeverExpires; x == nil || !x.Valid || x.Bool {
		t.Fatalf("got never expires %#v, want false", x)
	}
	if x := sub1.NumberOfBillingCycles; x == nil || !x.Valid || x.Int64 != 2 {
		t.Fatalf("got number billing cycles %#v, want 2", x)
	}
	if x := sub1.Price; x == nil || x.Scale != 2 || x.Unscaled != 100 {
		t.Fatalf("got price %#v, want 1.00", x)
	}
	if x := sub1.TrialPeriod; x == nil || !x.Valid || x.Bool {
		t.Fatalf("got trial period %#v, want false", x)
	}
	if x := sub1.Descriptor.Name; x != "Company Name*Product 1" {
		t.Fatalf("got descriptor name %#v, want Company Name*Product 1", x)
	}
	if x := sub1.Descriptor.Phone; x != "0000000000" {
		t.Fatalf("got descriptor phone %#v, want 0000000000", x)
	}
	if x := sub1.Descriptor.URL; x != "example.com" {
		t.Fatalf("got descriptor url %#v, want example.com", x)
	}

	// Update
	sub2, err := g.Update(&SubscriptionRequest{
		Id:     sub1.Id,
		PlanId: "test_plan_2",
		Options: &SubscriptionOptions{
			ProrateCharges:                       true,
			RevertSubscriptionOnProrationFailure: true,
			StartImmediately:                     true,
		},
	})

	t.Log("sub2", sub2)

	if err != nil {
		t.Fatal(err)
	}
	if sub2.Id != sub1.Id {
		t.Fatal(sub2.Id)
	}
	if x := sub2.PlanId; x != "test_plan_2" {
		t.Fatal(x)
	}

	// Find
	sub3, err := g.Find(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if sub3.Id != sub1.Id {
		t.Fatal(sub3.Id)
	}

	// Cancel
	_, err = g.Cancel(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSubscriptionAllFieldsWithBillingDayOfMonthNeverExpires(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(customer)

	g := testGateway.Subscription()

	// Create
	billingDayOfMonth := nullable.NewNullInt64(15, true)
	neverExpires := nullable.NewNullBool(true, true)
	sub1, err := g.Create(&SubscriptionRequest{
		PaymentMethodToken: paymentMethod.GetToken(),
		PlanId:             "test_plan",
		MerchantAccountId:  testMerchantAccountId,
		BillingDayOfMonth:  &billingDayOfMonth,
		NeverExpires:       &neverExpires,
		Price:              NewDecimal(100, 2),
		Descriptor: &Descriptor{
			Name:  "Company Name*Product 1",
			Phone: "0000000000",
			URL:   "example.com",
		},
	})

	t.Log("sub1", sub1)

	if err != nil {
		t.Fatal(err)
	}
	if sub1.Id == "" {
		t.Fatal("invalid subscription id")
	}
	if sub1.BillingDayOfMonth != "15" {
		t.Fatalf("got billing day of month %#v, want %#v", sub1.BillingDayOfMonth, "15")
	}
	if x := sub1.NeverExpires; x == nil || !x.Valid || !x.Bool {
		t.Fatalf("got never expires %#v, want true", x)
	}
	if x := sub1.NumberOfBillingCycles; x == nil || x.Valid {
		t.Fatalf("got number billing cycles %#v, didn't want", x)
	}
	if x := sub1.Price; x == nil || x.Scale != 2 || x.Unscaled != 100 {
		t.Fatalf("got price %#v, want 1.00", x)
	}
	if x := sub1.TrialPeriod; x == nil || !x.Valid || x.Bool {
		t.Fatalf("got trial period %#v, want false", x)
	}
	if x := sub1.Descriptor.Name; x != "Company Name*Product 1" {
		t.Fatalf("got descriptor name %#v, want Company Name*Product 1", x)
	}
	if x := sub1.Descriptor.Phone; x != "0000000000" {
		t.Fatalf("got descriptor phone %#v, want 0000000000", x)
	}
	if x := sub1.Descriptor.URL; x != "example.com" {
		t.Fatalf("got descriptor url %#v, want example.com", x)
	}

	// Update
	sub2, err := g.Update(&SubscriptionRequest{
		Id:     sub1.Id,
		PlanId: "test_plan_2",
		Options: &SubscriptionOptions{
			ProrateCharges:                       true,
			RevertSubscriptionOnProrationFailure: true,
			StartImmediately:                     true,
		},
	})

	t.Log("sub2", sub2)

	if err != nil {
		t.Fatal(err)
	}
	if sub2.Id != sub1.Id {
		t.Fatal(sub2.Id)
	}
	if x := sub2.PlanId; x != "test_plan_2" {
		t.Fatal(x)
	}

	// Find
	sub3, err := g.Find(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if sub3.Id != sub1.Id {
		t.Fatal(sub3.Id)
	}

	// Cancel
	_, err = g.Cancel(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSubscriptionAllFieldsWithFirstBillingDate(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(customer)

	g := testGateway.Subscription()

	// Create
	firstBillingDate := fmt.Sprintf("%d-12-31", time.Now().Year())
	numberOfBillingCycles := nullable.NewNullInt64(2, true)
	sub1, err := g.Create(&SubscriptionRequest{
		PaymentMethodToken:    paymentMethod.GetToken(),
		PlanId:                "test_plan",
		MerchantAccountId:     testMerchantAccountId,
		FirstBillingDate:      firstBillingDate,
		NumberOfBillingCycles: &numberOfBillingCycles,
		Price: NewDecimal(100, 2),
		Descriptor: &Descriptor{
			Name:  "Company Name*Product 1",
			Phone: "0000000000",
			URL:   "example.com",
		},
	})

	t.Log("sub1", sub1)

	if err != nil {
		t.Fatal(err)
	}
	if sub1.Id == "" {
		t.Fatal("invalid subscription id")
	}
	if sub1.BillingDayOfMonth != "31" {
		t.Fatalf("got billing day of month %#v, want %#v", sub1.BillingDayOfMonth, "31")
	}
	if sub1.FirstBillingDate != firstBillingDate {
		t.Fatalf("got first billing date %#v, want %#v", sub1.FirstBillingDate, firstBillingDate)
	}
	if x := sub1.NeverExpires; x == nil || !x.Valid || x.Bool {
		t.Fatalf("got never expires %#v, want false", x)
	}
	if x := sub1.NumberOfBillingCycles; x == nil || !x.Valid || x.Int64 != 2 {
		t.Fatalf("got number billing cycles %#v, want 2", x)
	}
	if x := sub1.Price; x == nil || x.Scale != 2 || x.Unscaled != 100 {
		t.Fatalf("got price %#v, want 1.00", x)
	}
	if x := sub1.TrialPeriod; x == nil || !x.Valid || x.Bool {
		t.Fatalf("got trial period %#v, want false", x)
	}
	if x := sub1.Descriptor.Name; x != "Company Name*Product 1" {
		t.Fatalf("got descriptor name %#v, want Company Name*Product 1", x)
	}
	if x := sub1.Descriptor.Phone; x != "0000000000" {
		t.Fatalf("got descriptor phone %#v, want 0000000000", x)
	}
	if x := sub1.Descriptor.URL; x != "example.com" {
		t.Fatalf("got descriptor url %#v, want example.com", x)
	}

	// Update
	sub2, err := g.Update(&SubscriptionRequest{
		Id:     sub1.Id,
		PlanId: "test_plan_2",
		Options: &SubscriptionOptions{
			ProrateCharges:                       true,
			RevertSubscriptionOnProrationFailure: true,
			StartImmediately:                     true,
		},
	})

	t.Log("sub2", sub2)

	if err != nil {
		t.Fatal(err)
	}
	if sub2.Id != sub1.Id {
		t.Fatal(sub2.Id)
	}
	if x := sub2.PlanId; x != "test_plan_2" {
		t.Fatal(x)
	}

	// Find
	sub3, err := g.Find(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if sub3.Id != sub1.Id {
		t.Fatal(sub3.Id)
	}

	// Cancel
	_, err = g.Cancel(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSubscriptionAllFieldsWithFirstBillingDateNeverExpires(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(customer)

	g := testGateway.Subscription()

	// Create
	firstBillingDate := fmt.Sprintf("%d-12-31", time.Now().Year())
	neverExpires := nullable.NewNullBool(true, true)
	sub1, err := g.Create(&SubscriptionRequest{
		PaymentMethodToken: paymentMethod.GetToken(),
		PlanId:             "test_plan",
		MerchantAccountId:  testMerchantAccountId,
		FirstBillingDate:   firstBillingDate,
		NeverExpires:       &neverExpires,
		Price:              NewDecimal(100, 2),
		Descriptor: &Descriptor{
			Name:  "Company Name*Product 1",
			Phone: "0000000000",
			URL:   "example.com",
		},
	})

	t.Log("sub1", sub1)

	if err != nil {
		t.Fatal(err)
	}
	if sub1.Id == "" {
		t.Fatal("invalid subscription id")
	}
	if sub1.BillingDayOfMonth != "31" {
		t.Fatalf("got billing day of month %#v, want %#v", sub1.BillingDayOfMonth, "31")
	}
	if sub1.FirstBillingDate != firstBillingDate {
		t.Fatalf("got first billing date %#v, want %#v", sub1.FirstBillingDate, firstBillingDate)
	}
	if x := sub1.NeverExpires; x == nil || !x.Valid || !x.Bool {
		t.Fatalf("got never expires %#v, want true", x)
	}
	if x := sub1.NumberOfBillingCycles; x == nil || x.Valid {
		t.Fatalf("got number billing cycles %#v, didn't want", x)
	}
	if x := sub1.Price; x == nil || x.Scale != 2 || x.Unscaled != 100 {
		t.Fatalf("got price %#v, want 1.00", x)
	}
	if x := sub1.TrialPeriod; x == nil || !x.Valid || x.Bool {
		t.Fatalf("got trial period %#v, want false", x)
	}
	if x := sub1.Descriptor.Name; x != "Company Name*Product 1" {
		t.Fatalf("got descriptor name %#v, want Company Name*Product 1", x)
	}
	if x := sub1.Descriptor.Phone; x != "0000000000" {
		t.Fatalf("got descriptor phone %#v, want 0000000000", x)
	}
	if x := sub1.Descriptor.URL; x != "example.com" {
		t.Fatalf("got descriptor url %#v, want example.com", x)
	}

	// Update
	sub2, err := g.Update(&SubscriptionRequest{
		Id:     sub1.Id,
		PlanId: "test_plan_2",
		Options: &SubscriptionOptions{
			ProrateCharges:                       true,
			RevertSubscriptionOnProrationFailure: true,
			StartImmediately:                     true,
		},
	})

	t.Log("sub2", sub2)

	if err != nil {
		t.Fatal(err)
	}
	if sub2.Id != sub1.Id {
		t.Fatal(sub2.Id)
	}
	if x := sub2.PlanId; x != "test_plan_2" {
		t.Fatal(x)
	}

	// Find
	sub3, err := g.Find(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if sub3.Id != sub1.Id {
		t.Fatal(sub3.Id)
	}

	// Cancel
	_, err = g.Cancel(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSubscriptionAllFieldsWithTrialPeriod(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(customer)

	g := testGateway.Subscription()

	// Create
	trialPeriod := nullable.NewNullBool(true, true)
	firstBillingDate := time.Now().AddDate(0, 0, 7)
	numberOfBillingCycles := nullable.NewNullInt64(2, true)
	sub1, err := g.Create(&SubscriptionRequest{
		PaymentMethodToken:    paymentMethod.GetToken(),
		PlanId:                "test_plan",
		MerchantAccountId:     testMerchantAccountId,
		TrialPeriod:           &trialPeriod,
		TrialDuration:         "7",
		TrialDurationUnit:     SubscriptionTrialDurationUnitDay,
		NumberOfBillingCycles: &numberOfBillingCycles,
		Price: NewDecimal(100, 2),
		Descriptor: &Descriptor{
			Name:  "Company Name*Product 1",
			Phone: "0000000000",
			URL:   "example.com",
		},
	})

	t.Log("sub1", sub1)

	if err != nil {
		t.Fatal(err)
	}
	if sub1.Id == "" {
		t.Fatal("invalid subscription id")
	}
	if sub1.BillingDayOfMonth != fmt.Sprintf("%d", firstBillingDate.Day()) {
		t.Fatalf("got billing day of month %#v, want %#v", sub1.BillingDayOfMonth, firstBillingDate.Day())
	}
	if sub1.FirstBillingDate != firstBillingDate.Format("2006-01-02") {
		t.Fatalf("got first billing date %#v, want %#v", sub1.FirstBillingDate, firstBillingDate)
	}
	if x := sub1.NeverExpires; x == nil || !x.Valid || x.Bool {
		t.Fatalf("got never expires %#v, want false", x)
	}
	if x := sub1.NumberOfBillingCycles; x == nil || !x.Valid || x.Int64 != 2 {
		t.Fatalf("got number billing cycles %#v, want 2", x)
	}
	if x := sub1.Price; x == nil || x.Scale != 2 || x.Unscaled != 100 {
		t.Fatalf("got price %#v, want 1.00", x)
	}
	if x := sub1.TrialPeriod; x == nil || !x.Valid || !x.Bool {
		t.Fatalf("got trial period %#v, want false", x)
	}
	if sub1.TrialDuration != "7" {
		t.Fatalf("got trial duration %#v, want 7", sub1.TrialDuration)
	}
	if sub1.TrialDurationUnit != SubscriptionTrialDurationUnitDay {
		t.Fatalf("got trial duration unit %#v, want day", sub1.TrialDurationUnit)
	}
	if x := sub1.Descriptor.Name; x != "Company Name*Product 1" {
		t.Fatalf("got descriptor name %#v, want Company Name*Product 1", x)
	}
	if x := sub1.Descriptor.Phone; x != "0000000000" {
		t.Fatalf("got descriptor phone %#v, want 0000000000", x)
	}
	if x := sub1.Descriptor.URL; x != "example.com" {
		t.Fatalf("got descriptor url %#v, want example.com", x)
	}

	// Update
	sub2, err := g.Update(&SubscriptionRequest{
		Id:     sub1.Id,
		PlanId: "test_plan_2",
		Options: &SubscriptionOptions{
			ProrateCharges:                       true,
			RevertSubscriptionOnProrationFailure: true,
			StartImmediately:                     true,
		},
	})

	t.Log("sub2", sub2)

	if err != nil {
		t.Fatal(err)
	}
	if sub2.Id != sub1.Id {
		t.Fatal(sub2.Id)
	}
	if x := sub2.PlanId; x != "test_plan_2" {
		t.Fatal(x)
	}

	// Find
	sub3, err := g.Find(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if sub3.Id != sub1.Id {
		t.Fatal(sub3.Id)
	}

	// Cancel
	_, err = g.Cancel(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSubscriptionAllFieldsWithTrialPeriodNeverExpires(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(customer)

	g := testGateway.Subscription()

	// Create
	trialPeriod := nullable.NewNullBool(true, true)
	firstBillingDate := time.Now().AddDate(0, 0, 7)
	neverExpires := nullable.NewNullBool(true, true)
	sub1, err := g.Create(&SubscriptionRequest{
		PaymentMethodToken: paymentMethod.GetToken(),
		PlanId:             "test_plan",
		MerchantAccountId:  testMerchantAccountId,
		TrialPeriod:        &trialPeriod,
		TrialDuration:      "7",
		TrialDurationUnit:  SubscriptionTrialDurationUnitDay,
		NeverExpires:       &neverExpires,
		Price:              NewDecimal(100, 2),
		Descriptor: &Descriptor{
			Name:  "Company Name*Product 1",
			Phone: "0000000000",
			URL:   "example.com",
		},
	})

	t.Log("sub1", sub1)

	if err != nil {
		t.Fatal(err)
	}
	if sub1.Id == "" {
		t.Fatal("invalid subscription id")
	}
	if sub1.BillingDayOfMonth != fmt.Sprintf("%d", firstBillingDate.Day()) {
		t.Fatalf("got billing day of month %#v, want %#v", sub1.BillingDayOfMonth, firstBillingDate.Day())
	}
	if sub1.FirstBillingDate != firstBillingDate.Format("2006-01-02") {
		t.Fatalf("got first billing date %#v, want %#v", sub1.FirstBillingDate, firstBillingDate)
	}
	if x := sub1.NeverExpires; x == nil || !x.Valid || !x.Bool {
		t.Fatalf("got never expires %#v, want true", x)
	}
	if x := sub1.NumberOfBillingCycles; x == nil || x.Valid {
		t.Fatalf("got number billing cycles %#v, didn't want", x)
	}
	if x := sub1.Price; x == nil || x.Scale != 2 || x.Unscaled != 100 {
		t.Fatalf("got price %#v, want 1.00", x)
	}
	if x := sub1.TrialPeriod; x == nil || !x.Valid || !x.Bool {
		t.Fatalf("got trial period %#v, want false", x)
	}
	if sub1.TrialDuration != "7" {
		t.Fatalf("got trial duration %#v, want 7", sub1.TrialDuration)
	}
	if sub1.TrialDurationUnit != SubscriptionTrialDurationUnitDay {
		t.Fatalf("got trial duration unit %#v, want day", sub1.TrialDurationUnit)
	}
	if x := sub1.Descriptor.Name; x != "Company Name*Product 1" {
		t.Fatalf("got descriptor name %#v, want Company Name*Product 1", x)
	}
	if x := sub1.Descriptor.Phone; x != "0000000000" {
		t.Fatalf("got descriptor phone %#v, want 0000000000", x)
	}
	if x := sub1.Descriptor.URL; x != "example.com" {
		t.Fatalf("got descriptor url %#v, want example.com", x)
	}

	// Update
	sub2, err := g.Update(&SubscriptionRequest{
		Id:     sub1.Id,
		PlanId: "test_plan_2",
		Options: &SubscriptionOptions{
			ProrateCharges:                       true,
			RevertSubscriptionOnProrationFailure: true,
			StartImmediately:                     true,
		},
	})

	t.Log("sub2", sub2)

	if err != nil {
		t.Fatal(err)
	}
	if sub2.Id != sub1.Id {
		t.Fatal(sub2.Id)
	}
	if x := sub2.PlanId; x != "test_plan_2" {
		t.Fatal(x)
	}

	// Find
	sub3, err := g.Find(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if sub3.Id != sub1.Id {
		t.Fatal(sub3.Id)
	}

	// Cancel
	_, err = g.Cancel(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
}
