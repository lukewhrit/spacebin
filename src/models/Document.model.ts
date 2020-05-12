import { Schema, model, Document, Model } from 'mongoose'
import config from '../config'

declare interface GlueDocumentInterface extends Document {
  content: string;
  key: string;
}

export type GlueDocumentModel = Model<GlueDocumentInterface>

export class GlueDocument {
  private _model: Model<GlueDocumentInterface>

  constructor () {
    const schema = new Schema({
      content: {
        type: String,
        required: true,
        min: 2,
        max: config.options.maxLength
      },
      key: {
        type: String,
        required: true,
        max: config.options.keyLength,
        min: config.options.keyLength
      }
    })

    this._model = model<GlueDocumentInterface>('GlueDocument', schema)
  }

  public get model (): Model<GlueDocumentInterface> {
    return this._model
  }
}
