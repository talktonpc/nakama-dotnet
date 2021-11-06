/**
* Copyright 2021 The Nakama Authors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

using Nakama;
using System.Collections.Generic;
using System.Linq;

namespace NakamaSync
{
    internal class HostMigrator : ISyncService
    {
        public SyncErrorHandler ErrorHandler { get; set; }
        public ILogger Logger { get; set; }

        private VarRegistry _registry;
        private readonly EnvelopeBuilder _builder;

        internal HostMigrator(VarRegistry registry, EnvelopeBuilder builder)
        {
            _registry = registry;
            _builder = builder;
        }

        public void Subscribe(PresenceTracker presenceTracker, HostTracker hostTracker)
        {
            hostTracker.OnHostChanged += (evt) =>
            {
                var self = presenceTracker.GetSelf();
                if (evt.OldHost != null && evt.NewHost?.UserId == self.UserId)
                {
                    // pick up where the old host left off in terms of validating values.
                    ValidatePendingVars(_registry);
                    UpdateVarHost(_registry, true);
                }
                else if (evt.OldHost?.UserId == self.UserId)
                {
                    UpdateVarHost(_registry, false);
                }
            };
        }

        private void ValidatePendingVars(VarRegistry registry)
        {
            ValidatePendingVars<bool>(registry.Bools.Values.SelectMany(l => l), env => env.SharedBoolAcks);
            ValidatePendingVars<float>(registry.Floats.Values.SelectMany(l => l), env => env.SharedFloatAcks);
            ValidatePendingVars<int>(registry.Ints.Values.SelectMany(l => l), env => env.SharedIntAcks);
            ValidatePendingVars<string>(registry.Strings.Values.SelectMany(l => l), env => env.SharedStringAcks);
            ValidatePendingVars<object>(registry.Objects.Values.SelectMany(l => l), env => env.SharedObjectAcks);

            _builder.SendEnvelope();
        }

        private void ValidatePendingVars<T>(IEnumerable<IVar<T>> vars, AckAccessor ackAccessor)
        {
            foreach (var var in vars)
            {
                _builder.AddAck(ackAccessor, var.Key);
            }
        }

        private void ValidatePendingVars<T>(Dictionary<string, List<IVar<T>>> vars, AckAccessor ackAccessor)
        {
            // TODO validate each var individually.
            foreach (var kvp in vars)
            {
                _builder.AddAck(ackAccessor, kvp.Key);
            }
        }

        private void UpdateVarHost(VarRegistry varRegistry, bool isHost)
        {
            UpdateVarHost(varRegistry.Bools.Values.SelectMany(l => l), isHost);
            UpdateVarHost(varRegistry.Floats.Values.SelectMany(l => l), isHost);
            UpdateVarHost(varRegistry.Ints.Values.SelectMany(l => l), isHost);
            UpdateVarHost(varRegistry.Strings.Values.SelectMany(l => l), isHost);
        }

        private void UpdateVarHost<T>(IEnumerable<IVar<T>> vars, bool isHost)
        {
            foreach (var var in vars)
            {
                var.IsHost = isHost;
            }
        }
    }
}
