export class UsersModel {
  user = {
    'model': {
      'username': 'username',
      'password': 'password',
      'role': 'user'
    },
    'schema': {
      'properties': {
        'username': {
          'type': 'string'
        },
        'password': {
          'type': 'string',
        },
        'role': {
          'type': 'string',
          'enum': ['user', 'admin'],
          'default': 'user'
        }
      }
    }
  };
}
