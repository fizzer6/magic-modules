package bigtable_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/bigtable"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccBigtableTable_basic(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable(instanceName, tableName),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableTable_splitKeys(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_splitKeys(instanceName, tableName),
			},
			{
				ResourceName:            "google_bigtable_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"split_keys"},
			},
		},
	})
}

func TestAccBigtableTable_family(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	family := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_family(instanceName, tableName, family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableTable_familyType(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	family := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_familyType(instanceName, tableName, family, "intmax"),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigtableTable_familyType(instanceName, tableName, family, `{
					"aggregateType": {
						"max": {},
						"inputType": {
							"int64Type": {
								"encoding": {
									"bigEndianBytes": {}
								}
							}
						}
					}
				}`),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigtableTable_familyType(instanceName, tableName, family, "intmax"),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccBigtableTable_familyType(instanceName, tableName, family, "intmin"),
				ExpectError: regexp.MustCompile("Immutable fields 'value_type.aggregate_type' cannot be updated"),
			},
		},
	})
}

func TestAccBigtableTable_testTableWithRowKeySchema(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	family := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_rowKeySchema(instanceName, tableName, family, `{
					"structType": {
						"fields": [{
							"fieldName": "myfield",
							"type": {
								"stringType": { "encoding": { "utf8Bytes": { } } }
							}
						}],
						"encoding": { "orderedCodeBytes": { } }
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableRowKeySchemaExists(t, "google_bigtable_table.table", true),
				),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// In-place modification is not accepted
				Config: testAccBigtableTable_rowKeySchema(instanceName, tableName, family, `{
					"structType": {
						"fields": [{
							"fieldName": "newfieldname",
							"type": {
								"stringType": { "encoding": { "utf8Bytes": { } } }
							}
						}],
						"encoding": { "orderedCodeBytes": { } }
					}
				}`),
				ExpectError: regexp.MustCompile(".*Row key schema in-place modification is not allowed.*"),
			},
			{
				// Removing the schema is ok
				Config: testAccBigtableTable_family(instanceName, tableName, family),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableRowKeySchemaExists(t, "google_bigtable_table.table", false),
				),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Set the schema to a new one is ok
			{
				Config: testAccBigtableTable_rowKeySchema(instanceName, tableName, family, `{
					"structType": {
						"fields": [
						    {
								"fieldName": "mystringfield",
								"type": {
									"stringType": { "encoding": { "utf8Bytes": { } } }
								}
							},
							{
								"fieldName": "myintfield",
								"type": {
									"int64Type": { "encoding": { "bigEndianBytes": { } } }
								}
							}
						],
						"encoding": { "delimitedBytes": { "delimiter": "Iw==" } }
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableRowKeySchemaExists(t, "google_bigtable_table.table", true),
				),
			},
		},
	})
}

func TestAccBigtableTable_deletion_protection_protected(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	family := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			// creating a table with a column family and deletion protection equals to protected
			{
				Config: testAccBigtableTable_deletion_protection(instanceName, tableName, "PROTECTED", family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// it is not possible to delete column families in the table with deletion protection equals to protected
			{
				Config:      testAccBigtableTable(instanceName, tableName),
				ExpectError: regexp.MustCompile(".*deletion protection field is set to true.*"),
			},
			// it is not possible to delete the table because of deletion protection equals to protected
			{
				Config:      testAccBigtableTable_destroyTable(instanceName),
				ExpectError: regexp.MustCompile(".*deletion protection field is set to true.*"),
			},
			// changing deletion protection field to unprotected without changing the column families
			// checking if the table and the column family exists
			{
				Config: testAccBigtableTable_deletion_protection(instanceName, tableName, "UNPROTECTED", family),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableColumnFamilyExists(t, "google_bigtable_table.table", family),
				),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// destroying the table is possible when deletion protection is equals to unprotected
			{
				Config: testAccBigtableTable_destroyTable(instanceName),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"},
			},
		},
	})
}

func TestAccBigtableTable_deletion_protection_unprotected(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	family := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			// creating a table with a column family and deletion protection equals to unprotected
			{
				Config: testAccBigtableTable_deletion_protection(instanceName, tableName, "UNPROTECTED", family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// removing the column family is possible because the deletion protection field is unprotected
			{
				Config: testAccBigtableTable(instanceName, tableName),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// changing the deletion protection field to protected
			{
				Config: testAccBigtableTable_deletion_protection(instanceName, tableName, "PROTECTED", family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// it is not possible to delete the table because of deletion protection equals to protected
			{
				Config:      testAccBigtableTable_destroyTable(instanceName),
				ExpectError: regexp.MustCompile(".*deletion protection field is set to true.*"),
			},
			// changing the deletion protection field to unprotected so that the sources can properly be destroyed
			{
				Config: testAccBigtableTable_deletion_protection(instanceName, tableName, "UNPROTECTED", family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableTable_change_stream_enable(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	family := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			// creating a table with a column family and change stream of 1 day
			{
				Config: testAccBigtableTable_change_stream_retention(instanceName, tableName, "24h0m0s", family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// it is not possible to delete the table because of change stream is enabled
			{
				Config:      testAccBigtableTable_destroyTable(instanceName),
				ExpectError: regexp.MustCompile(".*the change stream is enabled.*"),
			},
			// changing change stream retention value
			{
				Config: testAccBigtableTable_change_stream_retention(instanceName, tableName, "120h0m0s", family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// it is not possible to delete the table because of change stream is enabled
			{
				Config:      testAccBigtableTable_destroyTable(instanceName),
				ExpectError: regexp.MustCompile(".*the change stream is enabled.*"),
			},
			// disable changing change stream retention
			{
				Config: testAccBigtableTable_change_stream_retention(instanceName, tableName, "0", family),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableChangeStreamDisabled(t),
				),
			},
			// destroying the table is possible when change stream is disabled
			{
				Config: testAccBigtableTable_destroyTable(instanceName),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"},
			},
		},
	})
}

func TestAccBigtableTable_automated_backups(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	family := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			// Creating a table with automated backup disabled
			{
				Config: testAccBigtableTable_no_automated_backup_policy(instanceName, tableName, family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, false)),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with explicit disabling of automated backup
			{
				Config: testAccBigtableTable_automated_backups(instanceName, tableName, "0", "0", family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, false)),
			},
			{
				ResourceName:            "google_bigtable_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"automated_backup_policy"}, // ImportStateVerify doesn't use the CustomizeDiff function
			},
			// Update other table properties, leave automated backup policy untouched
			{
				Config: testAccBigtableTable_deletion_protection(instanceName, tableName, "PROTECTED", family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, false)),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update other table properties, leave automated backup policy untouched
			{
				Config: testAccBigtableTable_deletion_protection(instanceName, tableName, "UNPROTECTED", family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, false)),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Enable automated backup
			{
				Config: testAccBigtableTable_automated_backups(instanceName, tableName, "72h0m0s", "24h0m0s", family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, true)),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Disable automated backup
			{
				Config: testAccBigtableTable_automated_backups(instanceName, tableName, "0", "0", family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, false)),
			},
			{
				ResourceName:            "google_bigtable_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"automated_backup_policy"}, // ImportStateVerify doesn't use the CustomizeDiff function
			},
			// it is possible to delete the table when automated backup is disabled
			{
				Config: testAccBigtableTable_destroyTable(instanceName),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"},
			},
			// Creating a table with automated backups enabled
			{
				Config: testAccBigtableTable_automated_backups(instanceName, tableName, "72h0m0s", "24h0m0s", family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, true)),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Changing automated backup retention period value
			{
				Config: testAccBigtableTable_automated_backups(instanceName, tableName, "72h0m0s", "", family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, true)),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Changing automated backup frequency value
			{
				Config: testAccBigtableTable_automated_backups(instanceName, tableName, "", "24h0m0s", family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, true)),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Changing both automated backup retention period and frequency values
			{
				Config: testAccBigtableTable_automated_backups(instanceName, tableName, "72h0m0s", "24h0m0s", family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, true)),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Disabling automated backup explicitly
			{
				Config: testAccBigtableTable_automated_backups(instanceName, tableName, "0", "0", family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, false)),
			},
			{
				ResourceName:            "google_bigtable_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"automated_backup_policy"}, // ImportStateVerify doesn't use CustomizeDiff function
			},
			// Removing automated backup policy field has no effect (i.e. keeps automated backup disabled).
			{
				Config: testAccBigtableTable_no_automated_backup_policy(instanceName, tableName, family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, false)),
			},
			{
				ResourceName:            "google_bigtable_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"automated_backup_policy"}, // ImportStateVerify doesn't use CustomizeDiff function
			},
			// Renable automated backup
			{
				Config: testAccBigtableTable_automated_backups(instanceName, tableName, "72h0m0s", "24h0m0s", family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, true)),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Removing automated backup policy field has no effect (i.e. automated backup remains enabled)
			{
				Config: testAccBigtableTable_no_automated_backup_policy(instanceName, tableName, family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, true)),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Disabling automated backup
			{
				Config: testAccBigtableTable_automated_backups(instanceName, tableName, "0", "0", family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, false)),
			},
			{
				ResourceName:            "google_bigtable_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"automated_backup_policy"}, // ImportStateVerify doesn't use CustomizeDiff function
			},
			// Renable automated backup
			{
				Config: testAccBigtableTable_automated_backups(instanceName, tableName, "72h0m0s", "24h0m0s", family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, true)),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// it is possible to delete the table when automated backups is enabled
			{
				Config: testAccBigtableTable_destroyTable(instanceName),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"},
			},
		},
	})
}

func TestAccBigtableTable_automated_backups_explicitly_disabled_on_create(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	family := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			// Creating a table with automated backup explicitly disabled
			{
				Config: testAccBigtableTable_automated_backups(instanceName, tableName, "0", "0", family),
				Check:  resource.ComposeTestCheckFunc(verifyBigtableAutomatedBackupsEnablementState(t, false)),
			},
			{
				ResourceName:            "google_bigtable_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"automated_backup_policy"}, // ImportStateVerify doesn't use CustomizeDiff function
			},
			// it is possible to delete the table when automated backup is disabled
			{
				Config: testAccBigtableTable_destroyTable(instanceName),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"},
			},
		},
	})
}

func TestAccBigtableTable_familyMany(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	family := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_familyMany(instanceName, tableName, family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableTable_familyUpdate(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	family := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_familyMany(instanceName, tableName, family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigtableTable_familyUpdate(instanceName, tableName, family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckBigtableTableDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		var ctx = context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_bigtable_table" {
				continue
			}

			config := acctest.GoogleProviderConfig(t)
			c, err := config.BigTableClientFactory(config.UserAgent).NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
			if err != nil {
				// The instance is already gone
				return nil
			}

			_, err = c.TableInfo(ctx, rs.Primary.Attributes["name"])
			if err == nil {
				return fmt.Errorf("Table still present. Found %s in %s.", rs.Primary.Attributes["name"], rs.Primary.Attributes["instance_name"])
			}

			c.Close()
		}

		return nil
	}
}

func testAccBigtableColumnFamilyExists(t *testing.T, table_name_space, family string) resource.TestCheckFunc {
	var ctx = context.Background()
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[table_name_space]
		if !ok {
			return fmt.Errorf("Table not found: %s", table_name_space)
		}

		config := acctest.GoogleProviderConfig(t)
		c, err := config.BigTableClientFactory(config.UserAgent).NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
		if err != nil {
			return fmt.Errorf("Error starting admin client. %s", err)
		}

		defer c.Close()

		table, err := c.TableInfo(ctx, rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Error retrieving table. Could not find %s in %s.", rs.Primary.Attributes["name"], rs.Primary.Attributes["instance_name"])
		}
		families, err := bigtable.FlattenColumnFamily(table.FamilyInfos)
		if err != nil {
			return fmt.Errorf("Error flattening column families: %v", err)
		}
		for _, data := range families {
			if data["family"] != family {
				return fmt.Errorf("Error checking column family. Could not find column family %s in %s.", family, rs.Primary.Attributes["name"])
			}
		}

		return nil
	}
}

func testAccBigtableRowKeySchemaExists(t *testing.T, table_name_space string, expected_has_schema bool) resource.TestCheckFunc {
	ctx := context.Background()
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[table_name_space]
		if !ok {
			return fmt.Errorf("Table not found during schema check: %v", table_name_space)
		}

		config := acctest.GoogleProviderConfig(t)
		c, err := config.BigTableClientFactory(config.UserAgent).NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
		if err != nil {
			return fmt.Errorf("Error starting admin client %s", err)
		}
		defer c.Close()

		table, err := c.TableInfo(ctx, rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Error retrieving table. Could not find %s in %s", rs.Primary.Attributes["name"], rs.Primary.Attributes["instance_name"])
		}

		actual_has_schema := (table.RowKeySchema != nil)
		if actual_has_schema != expected_has_schema {
			return fmt.Errorf("expecting table to have row key schema to be %v, got %v", expected_has_schema, actual_has_schema)
		}

		return nil
	}
}

func testAccBigtableChangeStreamDisabled(t *testing.T) resource.TestCheckFunc {
	var ctx = context.Background()
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["google_bigtable_table.table"]
		if !ok {
			return fmt.Errorf("Table not found: %s", "google_bigtable_table.table")
		}

		config := acctest.GoogleProviderConfig(t)
		c, err := config.BigTableClientFactory(config.UserAgent).NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
		if err != nil {
			return fmt.Errorf("Error starting admin client. %s", err)
		}

		defer c.Close()

		table, err := c.TableInfo(ctx, rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Error retrieving table. Could not find %s in %s.", rs.Primary.Attributes["name"], rs.Primary.Attributes["instance_name"])
		}

		if table.ChangeStreamRetention != nil {
			return fmt.Errorf("Change Stream is expected to be disabled but it's not: %v", table)
		}

		return nil
	}
}

func verifyBigtableAutomatedBackupsEnablementState(t *testing.T, expectEnabled bool) resource.TestCheckFunc {
	var ctx = context.Background()
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["google_bigtable_table.table"]
		if !ok {
			return fmt.Errorf("Table not found: %s", "google_bigtable_table.table")
		}

		config := acctest.GoogleProviderConfig(t)
		c, err := config.BigTableClientFactory(config.UserAgent).NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
		if err != nil {
			return fmt.Errorf("Error starting admin client. %s", err)
		}

		defer c.Close()

		table, err := c.TableInfo(ctx, rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Error retrieving table. Could not find %s in %s.", rs.Primary.Attributes["name"], rs.Primary.Attributes["instance_name"])
		}
		if table.AutomatedBackupConfig != nil && !expectEnabled {
			return fmt.Errorf("Automated backup is expected to be disabled but it is not: %v", table)
		}
		if table.AutomatedBackupConfig == nil && expectEnabled {
			return fmt.Errorf("Automated backup is expected to be enabled but it is not: %v", table)
		}

		return nil
	}
}

func testAccBigtableTable(instanceName, tableName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  instance_type = "DEVELOPMENT"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id
}
`, instanceName, instanceName, tableName)
}

func testAccBigtableTable_splitKeys(instanceName, tableName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  instance_type = "DEVELOPMENT"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id
  split_keys    = ["a", "b", "c"]
}
`, instanceName, instanceName, tableName)
}

func testAccBigtableTable_family(instanceName, tableName, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name

  column_family {
    family = "%s"
  }
}
`, instanceName, instanceName, tableName, family)
}

func testAccBigtableTable_familyType(instanceName, tableName, family, familyType string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name

  column_family {
    family = "%s"
	type =  <<EOF
%s
EOF
  }
}
`, instanceName, instanceName, tableName, family, familyType)
}

func testAccBigtableTable_rowKeySchema(instanceName, tableName, family, rowKeySchema string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name

  column_family {
    family = "%s"
  }

  row_key_schema = <<EOF
%s 
EOF
}
`, instanceName, instanceName, tableName, family, rowKeySchema)
}

func testAccBigtableTable_deletion_protection(instanceName, tableName, deletionProtection, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name
  deletion_protection = "%s"

  column_family {
    family = "%s"
  }
}
`, instanceName, instanceName, tableName, deletionProtection, family)
}

func testAccBigtableTable_change_stream_retention(instanceName, tableName, changeStreamRetention, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name
  change_stream_retention = "%s"

  column_family {
    family = "%s"
  }
}
`, instanceName, instanceName, tableName, changeStreamRetention, family)
}

func testAccBigtableTable_automated_backups(instanceName, tableName, automatedBackupsRetentionPeriod, automatedBackupsFrequency, family string) string {
	var retentionPeriod string
	if automatedBackupsRetentionPeriod != "" {
		retentionPeriod = fmt.Sprintf(`retention_period = "%s"`, automatedBackupsRetentionPeriod)
	}
	var frequency string
	if automatedBackupsFrequency != "" {
		frequency = fmt.Sprintf(`frequency = "%s"`, automatedBackupsFrequency)
	}
	config := fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }
  instance_type = "DEVELOPMENT"
  deletion_protection = false
}
resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name
  automated_backup_policy {
	%s
	%s
  }
  column_family {
    family = "%s"
  }
}
`, instanceName, instanceName, tableName, retentionPeriod, frequency, family)
	return config
}

func testAccBigtableTable_no_automated_backup_policy(instanceName, tableName, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
	name = "%s"
	cluster {
	cluster_id = "%s"
	zone       = "us-central1-b"
	}
	instance_type = "DEVELOPMENT"
	deletion_protection = false
}
resource "google_bigtable_table" "table" {
	name          = "%s"
	instance_name = google_bigtable_instance.instance.name
	column_family {
	family = "%s"
	}
}
`, instanceName, instanceName, tableName, family)
}

func testAccBigtableTable_familyMany(instanceName, tableName, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name

  column_family {
    family = "%s-first"
  }

  column_family {
    family = "%s-second"
  }
}
`, instanceName, instanceName, tableName, family, family)
}

func testAccBigtableTable_familyUpdate(instanceName, tableName, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name

  column_family {
    family = "%s-third"
  }

  column_family {
    family = "%s-fourth"
  }

  column_family {
    family = "%s-second"
  }
}
`, instanceName, instanceName, tableName, family, family, family)
}

func testAccBigtableTable_destroyTable(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}
`, instanceName, instanceName)
}
