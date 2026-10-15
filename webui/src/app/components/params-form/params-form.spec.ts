import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ParamsFormComponent } from './params-form';
import { SchemaService } from '../../services/schema';
import { SomeJSONSchema } from 'ajv/dist/types/json-schema';

describe('ParamsForm', () => {
  let component: ParamsFormComponent;
  let fixture: ComponentFixture<ParamsFormComponent>;

  beforeEach(async () => {
    fixture = TestBed.createComponent(ParamsFormComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should create form for schema', async () => {
    const paramsSchema: SomeJSONSchema = {
      type: 'object',
      properties: {},
      required: [],
    };
    fixture.componentRef.setInput('paramsSchema', paramsSchema);
    await fixture.whenStable();
    expect(fixture.nativeElement.querySelector('form')).toBeTruthy();
  });

  it('form should have input for string property', async () => {
    const paramsSchema: SomeJSONSchema = {
      type: 'object',
      properties: {
        test: {
          type: 'string',
        },
      },
      required: [],
    };
    fixture.componentRef.setInput('paramsSchema', paramsSchema);
    await fixture.whenStable();
    const formEl = fixture.nativeElement.querySelector('form');
    expect(formEl).toBeDefined();
    const inputEl = formEl.querySelector('input');
    expect(inputEl).toBeDefined();
  });

  it('form should have select for string enum property', async () => {
    const paramsSchema: SomeJSONSchema = {
      type: 'object',
      properties: {
        test: {
          type: 'string',
          enum: ['one', 'two', 'three'],
        },
      },
      required: [],
    };
    fixture.componentRef.setInput('paramsSchema', paramsSchema);
    await fixture.whenStable();
    const formEl = fixture.nativeElement.querySelector('form');
    expect(formEl).toBeDefined();
    const selectEl = formEl.querySelector('select');
    expect(selectEl).toBeDefined();
  });

  it('form input should reflect data', async () => {
    const paramsSchema: SomeJSONSchema = {
      type: 'object',
      properties: {
        test: {
          type: 'string',
        },
      },
      required: [],
    };

    fixture.componentRef.setInput('paramsSchema', paramsSchema);
    await fixture.whenStable();
    fixture.componentRef.setInput('data', {
      test: 'hello',
    });
    await fixture.whenStable();
    const inputEl = fixture.nativeElement.querySelector('input');
    expect(inputEl).toBeDefined();
    expect(inputEl.value).toBe('hello');

    fixture.componentRef.setInput('data', {
      test: 'changed',
    });
    await fixture.whenStable();
    expect(inputEl.value).toBe('changed');
  });
});
