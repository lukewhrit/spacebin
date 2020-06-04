import { Entity, Column, PrimaryColumn } from 'typeorm'
import config from '../config'

@Entity()
export class Document {
  @PrimaryColumn({
    // Enforce key length
    length: config.options.keyLength
  })
  id: string

  @Column()
  content: string
}
