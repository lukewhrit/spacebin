import { Entity, Column, PrimaryGeneratedColumn } from 'typeorm'
import config from '../config'

@Entity()
export class Document {
  /*
   * This column provides an easy way to see the amount of documents stored.
   * It is non-essential, it could be removed.
   */
  @PrimaryGeneratedColumn()
  id: number

  @Column({
    // Enforce key length
    length: config.options.keyLength
  })
  key: string

  @Column()
  content: string
}
