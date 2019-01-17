import { storiesOf } from '@storybook/angular';

import { NewClusterComponent } from './new-cluster.component';
import { NEW_CLUSTER_MODULE_METADATA } from './new-cluster.component.metadata';

storiesOf('new cluster', module)
  .add('default', () => ({
    component: NewClusterComponent,
    moduleMetadata: NEW_CLUSTER_MODULE_METADATA
  }));
