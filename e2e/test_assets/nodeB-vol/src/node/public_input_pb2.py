# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: public_input.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x12public_input.proto\x12\x04main\"M\n\x07Testing\x12\x16\n\x0eis_early_reply\x18\x01 \x01(\x08\x12\x15\n\ris_early_exit\x18\x02 \x01(\x08\x12\x13\n\x0btest_stores\x18\x03 \x01(\x08\"-\n\x0eTestingResults\x12\x1b\n\x13test_stores_success\x18\x01 \x01(\x08\"7\n\x07Request\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x1e\n\x07testing\x18\x02 \x01(\x0b\x32\r.main.Testing\"o\n\x0cNodeBRequest\x12\x10\n\x08greeting\x18\x01 \x01(\t\x12\x1e\n\x07testing\x18\x02 \x01(\x0b\x32\r.main.Testing\x12-\n\x0ftesting_results\x18\x03 \x01(\x0b\x32\x14.main.TestingResults\"o\n\x0cNodeCRequest\x12\x10\n\x08greeting\x18\x01 \x01(\t\x12\x1e\n\x07testing\x18\x02 \x01(\x0b\x32\r.main.Testing\x12-\n\x0ftesting_results\x18\x03 \x01(\x0b\x32\x14.main.TestingResults\"k\n\x08Response\x12\x10\n\x08greeting\x18\x01 \x01(\t\x12\x1e\n\x07testing\x18\x02 \x01(\x0b\x32\r.main.Testing\x12-\n\x0ftesting_results\x18\x03 \x01(\x0b\x32\x14.main.TestingResults26\n\nEntrypoint\x12(\n\x05Greet\x12\r.main.Request\x1a\x0e.main.Response\"\x00\x42\tZ\x07./protob\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'public_input_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\007./proto'
  _TESTING._serialized_start=28
  _TESTING._serialized_end=105
  _TESTINGRESULTS._serialized_start=107
  _TESTINGRESULTS._serialized_end=152
  _REQUEST._serialized_start=154
  _REQUEST._serialized_end=209
  _NODEBREQUEST._serialized_start=211
  _NODEBREQUEST._serialized_end=322
  _NODECREQUEST._serialized_start=324
  _NODECREQUEST._serialized_end=435
  _RESPONSE._serialized_start=437
  _RESPONSE._serialized_end=544
  _ENTRYPOINT._serialized_start=546
  _ENTRYPOINT._serialized_end=600
# @@protoc_insertion_point(module_scope)
