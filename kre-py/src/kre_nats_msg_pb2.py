# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: kre_nats_msg.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='kre_nats_msg.proto',
  package='',
  syntax='proto3',
  serialized_options=_b('Z\005./kre'),
  serialized_pb=_b('\n\x12kre_nats_msg.proto\x1a\x19google/protobuf/any.proto\"\xf6\x01\n\x0eKreNatsMessage\x12\x13\n\x0btracking_id\x18\x01 \x01(\t\x12*\n\x08tracking\x18\x02 \x03(\x0b\x32\x18.KreNatsMessage.Tracking\x12\r\n\x05reply\x18\x03 \x01(\t\x12%\n\x07payload\x18\x04 \x01(\x0b\x32\x14.google.protobuf.Any\x12\r\n\x05\x65rror\x18\x05 \x01(\t\x12\x0f\n\x07replied\x18\x06 \x01(\x08\x12\x12\n\nearly_exit\x18\x07 \x01(\x08\x1a\x39\n\x08Tracking\x12\x11\n\tnode_name\x18\x01 \x01(\t\x12\r\n\x05start\x18\x02 \x01(\t\x12\x0b\n\x03\x65nd\x18\x03 \x01(\tB\x07Z\x05./kreb\x06proto3')
  ,
  dependencies=[google_dot_protobuf_dot_any__pb2.DESCRIPTOR,])




_KRENATSMESSAGE_TRACKING = _descriptor.Descriptor(
  name='Tracking',
  full_name='KreNatsMessage.Tracking',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='node_name', full_name='KreNatsMessage.Tracking.node_name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='start', full_name='KreNatsMessage.Tracking.start', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='end', full_name='KreNatsMessage.Tracking.end', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=239,
  serialized_end=296,
)

_KRENATSMESSAGE = _descriptor.Descriptor(
  name='KreNatsMessage',
  full_name='KreNatsMessage',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='tracking_id', full_name='KreNatsMessage.tracking_id', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='tracking', full_name='KreNatsMessage.tracking', index=1,
      number=2, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='reply', full_name='KreNatsMessage.reply', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='payload', full_name='KreNatsMessage.payload', index=3,
      number=4, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='error', full_name='KreNatsMessage.error', index=4,
      number=5, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='replied', full_name='KreNatsMessage.replied', index=5,
      number=6, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='early_exit', full_name='KreNatsMessage.early_exit', index=6,
      number=7, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_KRENATSMESSAGE_TRACKING, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=50,
  serialized_end=296,
)

_KRENATSMESSAGE_TRACKING.containing_type = _KRENATSMESSAGE
_KRENATSMESSAGE.fields_by_name['tracking'].message_type = _KRENATSMESSAGE_TRACKING
_KRENATSMESSAGE.fields_by_name['payload'].message_type = google_dot_protobuf_dot_any__pb2._ANY
DESCRIPTOR.message_types_by_name['KreNatsMessage'] = _KRENATSMESSAGE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

KreNatsMessage = _reflection.GeneratedProtocolMessageType('KreNatsMessage', (_message.Message,), dict(

  Tracking = _reflection.GeneratedProtocolMessageType('Tracking', (_message.Message,), dict(
    DESCRIPTOR = _KRENATSMESSAGE_TRACKING,
    __module__ = 'kre_nats_msg_pb2'
    # @@protoc_insertion_point(class_scope:KreNatsMessage.Tracking)
    ))
  ,
  DESCRIPTOR = _KRENATSMESSAGE,
  __module__ = 'kre_nats_msg_pb2'
  # @@protoc_insertion_point(class_scope:KreNatsMessage)
  ))
_sym_db.RegisterMessage(KreNatsMessage)
_sym_db.RegisterMessage(KreNatsMessage.Tracking)


DESCRIPTOR._options = None
# @@protoc_insertion_point(module_scope)
