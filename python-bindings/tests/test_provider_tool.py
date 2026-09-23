######################################################################
#
# File: python-bindings/tests/test_provider_tool.py
#
# Copyright 2026 Backblaze Inc. All Rights Reserved.
#
# License https://www.backblaze.com/using_b2_code.html
#
######################################################################

import unittest
from types import SimpleNamespace
from unittest.mock import Mock

from b2sdk.v3 import B2Api, InMemoryAccountInfo

from b2_terraform.provider_tool import Bucket


class BucketReadTest(unittest.TestCase):
    bucket_id = '5e3c3b5f8dcc43807ce50718'
    bucket_name = 'Nextcloud-Fastnetserv'
    account_id = 'ecbfdc30c578'

    def setUp(self):
        account_info = InMemoryAccountInfo()
        account_info.set_auth_data(
            account_id=self.account_id,
            auth_token='auth-token',
            api_url='https://api.example.invalid',
            download_url='https://download.example.invalid',
            recommended_part_size=100_000_000,
            absolute_minimum_part_size=5_000_000,
            application_key='application-key',
            realm='production',
            s3_api_url='https://s3.example.invalid',
            allowed={
                'buckets': [{'id': self.bucket_id, 'name': self.bucket_name}],
                'capabilities': [
                    'listBuckets',
                    'readBucketEncryption',
                    'readBucketRetentions',
                ],
                'namePrefix': None,
            },
            application_key_id='application-key-id',
        )
        self.api = B2Api(account_info=account_info)

        # b2sdk does this after authorizing a bucket-restricted application key.
        # It makes get_bucket_by_id/name return an ID/name-only Bucket object.
        self.api._populate_bucket_cache_from_key()
        self.api.session.list_buckets = Mock(
            return_value={'buckets': [self._bucket_api_response()]}
        )
        self.command = Bucket(SimpleNamespace(api=self.api))

    def _bucket_api_response(self):
        return {
            'accountId': self.account_id,
            'bucketId': self.bucket_id,
            'bucketName': self.bucket_name,
            'bucketType': 'allPrivate',
            'bucketInfo': {},
            'corsRules': [],
            'lifecycleRules': [],
            'options': ['s3'],
            'revision': 7,
            'defaultServerSideEncryption': {
                'isClientAuthorizedToRead': True,
                'value': {'algorithm': 'AES256', 'mode': 'SSE-B2'},
            },
            'fileLockConfiguration': {
                'isClientAuthorizedToRead': True,
                'value': {
                    'defaultRetention': {'mode': None, 'period': None},
                    'isFileLockEnabled': False,
                },
            },
        }

    def test_sdk_cached_bucket_is_incomplete(self):
        bucket = self.api.get_bucket_by_id(self.bucket_id)
        result = self.command._postprocess(bucket)

        self.assertIsNone(bucket.type_)
        self.assertEqual({'mode': None}, bucket.default_server_side_encryption.as_dict())
        self.assertEqual({'mode': 'unknown'}, bucket.default_retention.as_dict())
        self.assertNotIn('bucketType', result)
        self.assertEqual({'mode': None}, result['defaultServerSideEncryption'])
        self.assertEqual(
            {'defaultRetention': {'mode': 'unknown'}},
            result['fileLockConfiguration'],
        )
        with self.assertRaisesRegex(
            RuntimeError,
            'default_retention can only be set if is_file_lock_enabled is true',
        ):
            self.command._preprocess(
                file_lock_configuration=[
                    {
                        'is_file_lock_enabled': False,
                        'default_retention': [{'mode': 'unknown'}],
                    }
                ]
            )

    def test_resource_read_fetches_complete_bucket_metadata(self):
        result = self.command.resource_read(bucket_id=self.bucket_id)

        self.assertEqual('allPrivate', result['bucketType'])
        self.assertEqual(
            {'algorithm': 'AES256', 'mode': 'SSE-B2'},
            result['defaultServerSideEncryption'],
        )
        self.assertEqual(
            {'isFileLockEnabled': False},
            result['fileLockConfiguration'],
        )
        self.api.session.list_buckets.assert_called_once_with(
            self.account_id,
            bucket_name=None,
            bucket_id=self.bucket_id,
        )

    def test_data_source_read_fetches_complete_bucket_metadata(self):
        result = self.command.data_source_read(bucket_name=self.bucket_name)

        self.assertEqual('allPrivate', result['bucketType'])
        self.assertEqual(
            {'algorithm': 'AES256', 'mode': 'SSE-B2'},
            result['defaultServerSideEncryption'],
        )
        self.assertEqual(
            {'isFileLockEnabled': False},
            result['fileLockConfiguration'],
        )
        self.api.session.list_buckets.assert_called_once_with(
            self.account_id,
            bucket_name=self.bucket_name,
            bucket_id=None,
        )


if __name__ == '__main__':
    unittest.main()
