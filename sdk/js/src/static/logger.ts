/* eslint-disable */
// @ts-nocheck
// Code Generated with API SDK Generator. DO NOT EDIT.

type ConsoleMessage =
  | { type: 'group' | 'groupCollapsed'; value: unknown[]; groupId: string }
  | { type: 'groupEnd'; groupId: string }
  | { type: 'log'; value: unknown[]; groupId: string }
  | { type: 'time'; label: string; groupId: string }
  | { type: 'timeEnd'; label: string; duration: number; groupId: string }

class GroupedLogger {
  constructor(
    private groupId: string,
    private parent: Konsole,
  ) {}

  log(...message: unknown[]): void {
    this.parent.logToGroup(this.groupId, ...message)
  }

  time(label: string): void {
    this.parent.timeInGroup(this.groupId, label)
  }

  timeEnd(label: string): void {
    this.parent.timeEndInGroup(this.groupId, label)
  }

  group(...label: unknown[]): GroupedLogger {
    return this.parent.createNestedGroup(this.groupId, ...label)
  }

  groupCollapsed(...label: unknown[]): GroupedLogger {
    return this.parent.createNestedGroupCollapsed(this.groupId, ...label)
  }

  end(): void {
    this.parent.groupEnd(this.groupId)
  }
}

export class Konsole {
  private messageGroups: Map<string, ConsoleMessage[]> = new Map()
  private groupDepths: Map<string, number> = new Map()
  private timers: Map<string, Map<string, number>> = new Map() // Nested structure: groupId -> label -> startTime
  private groupHierarchy: Map<string, Set<string>> = new Map() // Track parent-child relationships

  private generateGroupId(): string {
    return `group_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`
  }

  private assertGroupExists(groupId: string): boolean {
    if (!this.messageGroups.has(groupId)) {
      console.warn(`Logger group ${groupId} not found - group may have been cleaned up`)
      return false
    }
    return true
  }

  private addToHierarchy(parentGroupId: string, childGroupId: string): void {
    let children = this.groupHierarchy.get(parentGroupId)
    if (!children) {
      children = new Set()
      this.groupHierarchy.set(parentGroupId, children)
    }
    children.add(childGroupId)
  }

  group(...label: unknown[]): GroupedLogger {
    const groupId = this.generateGroupId()
    this.messageGroups.set(groupId, [])
    this.groupDepths.set(groupId, 1)

    const messages = this.messageGroups.get(groupId)!
    messages.push({
      type: 'group',
      value: label,
      groupId,
    })

    return new GroupedLogger(groupId, this)
  }

  groupCollapsed(...label: unknown[]): GroupedLogger {
    const groupId = this.generateGroupId()
    this.messageGroups.set(groupId, [])
    this.groupDepths.set(groupId, 1)

    const messages = this.messageGroups.get(groupId)!
    messages.push({
      type: 'groupCollapsed',
      value: label,
      groupId,
    })

    return new GroupedLogger(groupId, this)
  }

  createNestedGroup(parentGroupId: string, ...label: unknown[]): GroupedLogger {
    if (!parentGroupId || typeof parentGroupId !== 'string' || parentGroupId.trim() === '') {
      console.error('Parent group ID must be a non-empty string')
      return this.group(...label)
    }

    if (!this.assertGroupExists(parentGroupId)) {
      // Parent group doesn't exist, create a new top-level group instead
      return this.group(...label)
    }
    const depth = this.groupDepths.get(parentGroupId)!

    const groupId = this.generateGroupId()
    this.messageGroups.set(groupId, [])
    this.groupDepths.set(groupId, depth + 1)
    this.addToHierarchy(parentGroupId, groupId)

    const messages = this.messageGroups.get(parentGroupId)!
    messages.push({
      type: 'group',
      value: label,
      groupId,
    })

    return new GroupedLogger(groupId, this)
  }

  createNestedGroupCollapsed(parentGroupId: string, ...label: unknown[]): GroupedLogger {
    if (!parentGroupId || typeof parentGroupId !== 'string' || parentGroupId.trim() === '') {
      console.error('Parent group ID must be a non-empty string')
      return this.groupCollapsed(...label)
    }

    if (!this.assertGroupExists(parentGroupId)) {
      // Parent group doesn't exist, create a new top-level group instead
      return this.groupCollapsed(...label)
    }
    const depth = this.groupDepths.get(parentGroupId)!

    const groupId = this.generateGroupId()
    this.messageGroups.set(groupId, [])
    this.groupDepths.set(groupId, depth + 1)
    this.addToHierarchy(parentGroupId, groupId)

    const messages = this.messageGroups.get(parentGroupId)!
    messages.push({
      type: 'groupCollapsed',
      value: label,
      groupId,
    })

    return new GroupedLogger(groupId, this)
  }

  logToGroup(groupId: string, ...message: unknown[]): void {
    if (!groupId || typeof groupId !== 'string' || groupId.trim() === '') {
      console.error('Group ID must be a non-empty string')
      console.log(...message)
      return
    }

    if (!this.assertGroupExists(groupId)) {
      // Group doesn't exist, log directly to console instead
      console.log(...message)
      return
    }
    const messages = this.messageGroups.get(groupId)!

    messages.push({
      type: 'log',
      value: message,
      groupId,
    })
  }

  timeInGroup(groupId: string, label: string): void {
    if (!label || typeof label !== 'string' || label.trim() === '') {
      console.error('Timer label must be a non-empty string')
      return
    }

    if (!this.assertGroupExists(groupId)) {
      // Group doesn't exist, start timer directly on console instead
      console.time(label)
      return
    }
    const messages = this.messageGroups.get(groupId)!

    let groupTimers = this.timers.get(groupId)
    if (!groupTimers) {
      groupTimers = new Map()
      this.timers.set(groupId, groupTimers)
    }

    if (groupTimers.has(label)) {
      console.error(`Timer "${label}" already exists in group ${groupId}`)
      return
    }

    groupTimers.set(label, performance.now())
    // Don't call console.time() here - store the info for later
    messages.push({
      type: 'time', // Changed type to distinguish it from actual console.time call
      label,
      groupId,
    })
  }

  timeEndInGroup(groupId: string, label: string): void {
    if (!label || typeof label !== 'string' || label.trim() === '') {
      console.error('Timer label must be a non-empty string')
      return
    }

    if (!this.assertGroupExists(groupId)) {
      // Group doesn't exist, end timer directly on console instead
      console.timeEnd(label)
      return
    }
    const messages = this.messageGroups.get(groupId)!

    const groupTimers = this.timers.get(groupId)
    if (!groupTimers) {
      console.error(`Timer "${label}" not found in group ${groupId}`)
      return
    }

    const startTime = groupTimers.get(label)
    if (startTime === undefined) {
      console.error(`Timer "${label}" not found in group ${groupId}`)
      return
    }

    const endTime = performance.now()
    const duration = endTime - startTime
    groupTimers.delete(label)

    // Store the calculated duration in the messages
    messages.push({
      type: 'timeEnd',
      label,
      duration,
      groupId,
    })
  }

  private cleanupGroup(groupId: string): void {
    // Cleanup the group itself
    this.messageGroups.delete(groupId)
    this.groupDepths.delete(groupId)

    // Cleanup timers for this group - now O(1) instead of O(n)
    this.timers.delete(groupId)

    // Cleanup child groups recursively
    const children = this.groupHierarchy.get(groupId)
    if (children) {
      children.forEach((childId) => {
        this.cleanupGroup(childId)
      })
      this.groupHierarchy.delete(groupId)
    }
  }

  groupEnd(groupId: string): void {
    if (!groupId || typeof groupId !== 'string' || groupId.trim() === '') {
      console.error('Group ID must be a non-empty string')
      return
    }

    if (!this.assertGroupExists(groupId)) {
      // Group doesn't exist, it may have already been cleaned up
      return
    }
    const messages = this.messageGroups.get(groupId)!

    // Process messages and output them with proper grouping
    this.processMessages(messages)

    // Cleanup this group and all its children
    this.cleanupGroup(groupId)
  }

  private processMessages(messages: ConsoleMessage[]): void {
    let groupStarted = false

    messages.forEach((message) => {
      switch (message.type) {
        case 'group':
          console.group(...message.value)
          groupStarted = true
          break
        case 'groupCollapsed':
          console.groupCollapsed(...message.value)
          groupStarted = true
          break
        case 'log':
          console.log(...message.value)
          break
        case 'time':
          console.log(`Timer '${message.label}' started`)
          break
        case 'timeEnd':
          console.log(`Timer '${message.label}': ${message.duration.toFixed(2)}ms`)
          break
        case 'groupEnd':
          if (groupStarted) {
            console.groupEnd()
            groupStarted = false
          }
          break
      }
    })

    // Ensure the group is closed if it was started but no explicit groupEnd was found
    if (groupStarted) {
      console.groupEnd()
    }
  }
}
