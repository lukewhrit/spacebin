/*
 * Copyright (C) 2020 The Spacebin Authors: notably Luke Whrit, Jack Dorland

 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.

 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import { Entity, Column, PrimaryColumn, CreateDateColumn } from 'typeorm'

@Entity()
export class Document {
  @PrimaryColumn()
  id: string

  @Column()
  content: string

  @CreateDateColumn()
  dateCreated: Date

  @Column({
    default: 'text'
  })
  extension: string
}
