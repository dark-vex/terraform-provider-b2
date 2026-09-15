//####################################################################
//
// File: b2/resource_b2_bucket_test.go
//
// Copyright 2021 Backblaze Inc. All Rights Reserved.
//
// License https://www.backblaze.com/using_b2_code.html
//
//####################################################################

package b2

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceB2Bucket_basic(t *testing.T) {
	resourceName := "b2_bucket.test"

	bucketName := acctest.RandomWithPrefix("test-b2-tfp")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceB2BucketConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "account_id", regexp.MustCompile("^[a-zA-Z0-9]{12}$")),
					resource.TestMatchResourceAttr(resourceName, "bucket_id", regexp.MustCompile("^[a-zA-Z0-9]{24}$")),
					resource.TestCheckResourceAttr(resourceName, "bucket_info.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(resourceName, "bucket_type", "allPublic"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "file_lock_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "file_lock_configuration.0.is_file_lock_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_server_side_encryption.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_server_side_encryption.0.mode", "SSE-B2"),
					resource.TestCheckResourceAttr(resourceName, "default_server_side_encryption.0.algorithm", "AES256"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "options.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "options.0", "s3"),
					resource.TestCheckResourceAttr(resourceName, "revision", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceB2Bucket_all(t *testing.T) {
	resourceName := "b2_bucket.test"

	bucketName := acctest.RandomWithPrefix("test-b2-tfp")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceB2BucketConfig_all(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "account_id", regexp.MustCompile("^[a-zA-Z0-9]{12}$")),
					resource.TestMatchResourceAttr(resourceName, "bucket_id", regexp.MustCompile("^[a-zA-Z0-9]{24}$")),
					resource.TestCheckResourceAttr(resourceName, "bucket_info.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "bucket_info.description", "the bucket"),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(resourceName, "bucket_type", "allPrivate"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.cors_rule_name", "downloadFromAnyOrigin"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_origins.0", "https"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_operations.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_operations.0", "b2_download_file_by_id"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_operations.1", "b2_download_file_by_name"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_operations.2", "s3_put"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_operations.3", "s3_head"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_operations.4", "s3_get"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.expose_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.expose_headers.0", "x-bz-content-sha1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_headers.0", "range"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.max_age_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceName, "file_lock_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "file_lock_configuration.0.is_file_lock_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "file_lock_configuration.0.default_retention.0.mode", "governance"),
					resource.TestCheckResourceAttr(resourceName, "file_lock_configuration.0.default_retention.0.period.0.duration", "7"),
					resource.TestCheckResourceAttr(resourceName, "file_lock_configuration.0.default_retention.0.period.0.unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "default_server_side_encryption.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_server_side_encryption.0.mode", "SSE-B2"),
					resource.TestCheckResourceAttr(resourceName, "default_server_side_encryption.0.algorithm", "AES256"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.0.file_name_prefix", "c/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.0.days_from_hiding_to_deleting", "1"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.0.days_from_uploading_to_hiding", "2"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.0.days_from_starting_to_canceling_unfinished_large_files", "3"),
					resource.TestCheckResourceAttr(resourceName, "options.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "options.0", "s3"),
					resource.TestCheckResourceAttr(resourceName, "revision", "3"),
				),
			},
		},
	})
}

func TestAccResourceB2Bucket_lifecycleRulesDefaults(t *testing.T) {
	resourceName := "b2_bucket.test"

	bucketName := acctest.RandomWithPrefix("test-b2-tfp")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceB2BucketConfig_lifecycleRulesDefaults(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "account_id", regexp.MustCompile("^[a-zA-Z0-9]{12}$")),
					resource.TestCheckResourceAttr(resourceName, "bucket_info.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(resourceName, "bucket_type", "allPublic"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.0.file_name_prefix", "a/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.0.days_from_hiding_to_deleting", "2"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.0.days_from_uploading_to_hiding", "0"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.0.days_from_starting_to_canceling_unfinished_large_files", "0"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.1.file_name_prefix", "b/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.1.days_from_hiding_to_deleting", "0"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.1.days_from_uploading_to_hiding", "2"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.1.days_from_starting_to_canceling_unfinished_large_files", "0"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.2.file_name_prefix", "c/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.2.days_from_hiding_to_deleting", "0"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.2.days_from_uploading_to_hiding", "0"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.2.days_from_starting_to_canceling_unfinished_large_files", "2"),
					resource.TestCheckResourceAttr(resourceName, "options.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "options.0", "s3"),
					resource.TestCheckResourceAttr(resourceName, "revision", "2"),
				),
			},
		},
	})
}

func TestAccResourceB2Bucket_update(t *testing.T) {
	resourceName := "b2_bucket.test"

	bucketName := acctest.RandomWithPrefix("test-b2-tfp")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceB2BucketConfig_basic(bucketName),
			},
			// We're testing a minimal change here to check if omitted optional
			// fields do not cause an update to fail.
			{
				Config: testAccResourceB2BucketConfig_basicWithFileInfo(bucketName),
			},
			{
				Config: testAccResourceB2BucketConfig_all(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "account_id", regexp.MustCompile("^[a-zA-Z0-9]{12}$")),
					resource.TestMatchResourceAttr(resourceName, "bucket_id", regexp.MustCompile("^[a-zA-Z0-9]{24}$")),
					resource.TestCheckResourceAttr(resourceName, "bucket_info.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "bucket_info.description", "the bucket"),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(resourceName, "bucket_type", "allPrivate"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.cors_rule_name", "downloadFromAnyOrigin"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_origins.0", "https"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_operations.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_operations.0", "b2_download_file_by_id"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_operations.1", "b2_download_file_by_name"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_operations.2", "s3_put"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_operations.3", "s3_head"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_operations.4", "s3_get"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.expose_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.expose_headers.0", "x-bz-content-sha1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.allowed_headers.0", "range"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.0.max_age_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceName, "default_server_side_encryption.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_server_side_encryption.0.mode", "SSE-B2"),
					resource.TestCheckResourceAttr(resourceName, "default_server_side_encryption.0.algorithm", "AES256"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.0.file_name_prefix", "c/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.0.days_from_hiding_to_deleting", "1"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.0.days_from_uploading_to_hiding", "2"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.0.days_from_starting_to_canceling_unfinished_large_files", "3"),
					resource.TestCheckResourceAttr(resourceName, "options.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "options.0", "s3"),
				),
			},
			{
				Config: testAccResourceB2BucketConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "account_id", regexp.MustCompile("^[a-zA-Z0-9]{12}$")),
					resource.TestMatchResourceAttr(resourceName, "bucket_id", regexp.MustCompile("^[a-zA-Z0-9]{24}$")),
					resource.TestCheckResourceAttr(resourceName, "bucket_info.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(resourceName, "bucket_type", "allPublic"),
					resource.TestCheckResourceAttr(resourceName, "cors_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "default_server_side_encryption.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_server_side_encryption.0.mode", "SSE-B2"),
					resource.TestCheckResourceAttr(resourceName, "default_server_side_encryption.0.algorithm", "AES256"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "options.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "options.0", "s3"),
					resource.TestCheckResourceAttr(resourceName, "revision", "2"),
				),
			},
		},
	})
}

func TestAccResourceB2Bucket_defaultRetention(t *testing.T) {
	resourceName := "b2_bucket.test"

	bucketName := acctest.RandomWithPrefix("test-b2-tfp")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceB2BucketConfig_defaultRetention(bucketName, false),
				ExpectError: regexp.MustCompile("default_retention"),
			},
			{
				Config: testAccResourceB2BucketConfig_defaultRetention(bucketName, true),
				Check:  resource.TestCheckResourceAttr(resourceName, "bucket_name", bucketName),
			},
		},
	})
}

func testAccResourceB2BucketConfig_defaultRetention(bucketName string, fileLockEnabled bool) string {
	return fmt.Sprintf(`
resource "b2_bucket" "test" {
	bucket_name = "%s"
	bucket_type = "allPublic"
	file_lock_configuration {
		is_file_lock_enabled = %t
		default_retention {
			mode = "governance"
			period {
				duration = 7
				unit     = "days"
			}
		}
	}
}
`, bucketName, fileLockEnabled)
}

func TestAccResourceB2Bucket_revisionUpdatedOnChange(t *testing.T) {
	resourceName := "b2_bucket.test"
	bucketName := acctest.RandomWithPrefix("test-b2-tfp")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceB2BucketConfig_basic(bucketName),
				Check:  resource.TestCheckResourceAttr(resourceName, "revision", "2"),
			},
			{
				Config: testAccResourceB2BucketConfig_basicPrivate(bucketName),
				Check:  resource.TestCheckResourceAttr(resourceName, "revision", "3"),
			},
		},
	})
}

func TestAccResourceB2Bucket_bucketInfoKeyCase(t *testing.T) {
	resourceName := "b2_bucket.test"
	bucketName := acctest.RandomWithPrefix("test-b2-tfp")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				// B2 stores the key in lower case, which must not show up as a permanent diff
				Config: testAccResourceB2BucketConfig_bucketInfo(bucketName, "ManagedBy"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket_info.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "bucket_info.managedby", "Terraform"),
				),
			},
			{
				// The key stored in lower case must not be planned as a change back to the config one
				Config:   testAccResourceB2BucketConfig_bucketInfo(bucketName, "ManagedBy"),
				PlanOnly: true,
			},
			{
				// Changing the key case only is not a change for B2
				Config:   testAccResourceB2BucketConfig_bucketInfo(bucketName, "managedby"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceB2Bucket_lifecycleRulesOrder(t *testing.T) {
	resourceName := "b2_bucket.test"
	bucketName := acctest.RandomWithPrefix("test-b2-tfp")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				// B2 does not return lifecycle rules in the order they were set,
				// the state must be stored in the config order regardless
				Config: testAccResourceB2BucketConfig_lifecycleRulesOrder(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.0.file_name_prefix", "uploads/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.0.days_from_uploading_to_hiding", "1"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.1.file_name_prefix", "files/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.1.days_from_hiding_to_deleting", "1"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.2.file_name_prefix", "thumbnails/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.2.days_from_uploading_to_hiding", "2"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.3.file_name_prefix", "avatars/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rules.3.days_from_hiding_to_deleting", "2"),
				),
			},
			{
				// The rule order must not be planned as a change
				Config:   testAccResourceB2BucketConfig_lifecycleRulesOrder(bucketName),
				PlanOnly: true,
			},
		},
	})
}

func testAccResourceB2BucketConfig_lifecycleRulesOrder(bucketName string) string {
	return fmt.Sprintf(`
resource "b2_bucket" "test" {
  bucket_name = "%s"
  bucket_type = "allPublic"

  lifecycle_rules {
    file_name_prefix = "uploads/"
    days_from_uploading_to_hiding = 1
  }
  lifecycle_rules {
    file_name_prefix = "files/"
    days_from_hiding_to_deleting = 1
  }
  lifecycle_rules {
    file_name_prefix = "thumbnails/"
    days_from_uploading_to_hiding = 2
  }
  lifecycle_rules {
    file_name_prefix = "avatars/"
    days_from_hiding_to_deleting = 2
  }
}
`, bucketName)
}

func testAccResourceB2BucketConfig_bucketInfo(bucketName string, key string) string {
	return fmt.Sprintf(`
resource "b2_bucket" "test" {
  bucket_name = "%s"
  bucket_type = "allPublic"
  bucket_info = {
    %s = "Terraform"
  }
}
`, bucketName, key)
}

func testAccResourceB2BucketConfig_basicPrivate(bucketName string) string {
	return fmt.Sprintf(`
resource "b2_bucket" "test" {
  bucket_name = "%s"
  bucket_type = "allPrivate"
}
`, bucketName)
}

func testAccResourceB2BucketConfig_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "b2_bucket" "test" {
  bucket_name = "%s"
  bucket_type = "allPublic"
}
`, bucketName)
}

func testAccResourceB2BucketConfig_basicWithFileInfo(bucketName string) string {
	return fmt.Sprintf(`
resource "b2_bucket" "test" {
  bucket_name = "%s"
  bucket_type = "allPublic"
  bucket_info = {
	key = "value"
  }
}
`, bucketName)
}

func testAccResourceB2BucketConfig_lifecycleRulesDefaults(bucketName string) string {
	return fmt.Sprintf(`
resource "b2_bucket" "test" {
  bucket_name = "%s"
  bucket_type = "allPublic"

  lifecycle_rules {
    file_name_prefix = "a/"
    days_from_hiding_to_deleting = 2
  }
  lifecycle_rules {
    file_name_prefix = "b/"
    days_from_uploading_to_hiding = 2
  }
  lifecycle_rules {
    file_name_prefix = "c/"
    days_from_starting_to_canceling_unfinished_large_files = 2
  }
}
`, bucketName)
}

func testAccResourceB2BucketConfig_all(bucketName string) string {
	return fmt.Sprintf(`
resource "b2_bucket" "test" {
  bucket_name = "%s"
  bucket_type = "allPrivate"
  bucket_info = {
    description = "the bucket"
  }
  cors_rules {
    cors_rule_name = "downloadFromAnyOrigin"
    allowed_origins = [
      "https"
    ]
    allowed_operations = [
      "b2_download_file_by_id",
      "b2_download_file_by_name",
	  "s3_put",
	  "s3_head",
	  "s3_get",
    ]
    expose_headers = ["x-bz-content-sha1"]
    allowed_headers = ["range"]
    max_age_seconds = 3600
  }
  file_lock_configuration {
	is_file_lock_enabled = true
	default_retention {
	  mode = "governance"
	  period {
		duration = 7
		unit = "days"
	  }
	}
  }
  default_server_side_encryption {
    mode = "SSE-B2"
    algorithm = "AES256"
  }
  lifecycle_rules {
    file_name_prefix = "c/"
    days_from_hiding_to_deleting = 1
    days_from_uploading_to_hiding = 2
    days_from_starting_to_canceling_unfinished_large_files = 3
  }
}
`, bucketName)
}
