# Mermaid Agent

A BER agent for creating and modifying various types of Mermaid diagrams.

## Features

- Flowchart diagram creation
- Sequence diagram generation
- Pie chart visualization
- Automated syntax validation
- Markdown-compatible output

## Prerequisites

- No additional configuration required
- Works out of the box with BER installation

## Usage

```
# For TUI mode
go run . tui

# For the GitHub Adapter
go run . webhook --debug
```

The agent provides diagramming capabilities through the following skills:
```go
FlowchartDiagram:
  - Node and connection creation
  - Direction control (TD, LR)
  - Shape and style customization
  - Relationship arrows

SequenceDiagram:
  - Actor definition
  - Message flow visualization
  - Activation blocks
  - Alternative paths
  - Loop sequences

PieChart:
  - Section definition
  - Value/percentage allocation
  - Title and labels
  - Data validation
```

## Validation

The agent performs several validations:
- Mermaid syntax correctness
- Diagram type appropriateness
- Label clarity and readability
- Relationship logic
