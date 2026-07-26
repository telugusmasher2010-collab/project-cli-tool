import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

/// Displays the list of expenses for the current group.
class ExpensePage extends StatefulWidget {
  const ExpensePage({super.key});

  @override
  State<ExpensePage> createState() => _ExpensePageState();
}

class _ExpensePageState extends State<ExpensePage> {
  // TODO: Replace with real data from Supabase.
  final List<Map<String, dynamic>> _expenses = [];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Expenses'),
        actions: [
          IconButton(
            icon: const Icon(Icons.group_add),
            onPressed: () {
              // TODO: Navigate to group management.
            },
          ),
        ],
      ),
      body: _expenses.isEmpty
          ? const Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.receipt_long, size: 64, color: Colors.grey),
                  SizedBox(height: 16),
                  Text(
                    'No expenses yet',
                    style: TextStyle(fontSize: 18, color: Colors.grey),
                  ),
                  SizedBox(height: 8),
                  Text(
                    'Tap + to add your first expense',
                    style: TextStyle(color: Colors.grey),
                  ),
                ],
              ),
            )
          : ListView.builder(
              itemCount: _expenses.length,
              itemBuilder: (context, index) {
                final expense = _expenses[index];
                return ListTile(
                  leading: const CircleAvatar(child: Icon(Icons.receipt)),
                  title: Text(expense['title'] ?? ''),
                  subtitle: Text('₹${expense['amount'] ?? 0}'),
                  trailing: Text(expense['date'] ?? ''),
                );
              },
            ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push('/add'),
        child: const Icon(Icons.add),
      ),
    );
  }
}
