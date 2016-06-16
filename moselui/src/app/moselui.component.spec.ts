import {
  beforeEachProviders,
  describe,
  expect,
  it,
  inject
} from '@angular/core/testing';
import { MoseluiAppComponent } from '../app/moselui.component';

beforeEachProviders(() => [MoseluiAppComponent]);

describe('App: Moselui', () => {
  it('should create the app',
      inject([MoseluiAppComponent], (app: MoseluiAppComponent) => {
    expect(app).toBeTruthy();
  }));
  
});
