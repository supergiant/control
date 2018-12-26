import { storiesOf } from '@storybook/angular';
import { FooterComponent } from './footer.component';

// example story
storiesOf('Footer', module)
  .add('Default', () => ({
    component: FooterComponent,
  }));
