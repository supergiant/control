export class UsersModel {
  user = {
    'model': {
      'username': '',
      'password': '',
      'role': 'user'
    },
    'schema': {
      'properties': {
        'username': {
          'description': 'User Name',
          'type': 'string'
        },
        'password': {
          'description': 'Password',
          'type': 'string',
          'widget': 'password'
        },
        'role': {
          'description': 'User Role',
          'type': 'string',
          'widget': 'select',
          'oneOf': [{
            'description': 'Admin', 'enum': ['admin']
          }, {
            'description': 'User', 'enum': ['user']
          }],
          'default': 'user'
        }
      }
    }
  };
}
