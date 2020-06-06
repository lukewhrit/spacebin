import { Entity, Column, PrimaryColumn } from 'typeorm'

@Entity()
export class Document {
  @PrimaryColumn()
  id: string

  @Column()
  content: string
}
